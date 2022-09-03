package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	// pretty is a flag to print the json output in a pretty format.
	pretty = flag.Bool("p", false, "pretty print")

	// file is a flag to specify a single file to parse.
	file = flag.String("f", "", "file")

	// help is a flag to print the help text.
	help = flag.Bool("h", false, "help")
)

var (
	// errPathIssue is the error message when the user provides both a file and a directory.
	errPathIssue = errors.New("ERROR: cannot specify both a file and a directory")

	// errNotAFile is the error message when the user provides a directory as a file.
	errNotAFile = errors.New("ERROR: required file, provided directory")

	// errUnknown is the error message when the error is not an expected type.
	errUnknown = errors.New("ERROR: unknown error")
)

// Detail is a struct that contains the details of a test.
type Detail struct {
	Name         string    `json:"name"`
	FileName     string    `json:"fileName"`
	RelativePath string    `json:"relativePath"`
	AbsolutePath string    `json:"absolutePath"`
	Line         int       `json:"line"`
	Pos          token.Pos `json:"pos"`
}

func main() {
	flag.Parse()

	err := Process(&args{
		file:   *file,
		dirs:   flag.Args(),
		help:   *help,
		pretty: *pretty,
	}, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// args is a struct that contains the arguments provided by the user.
type args struct {
	file   string
	dirs   []string
	help   bool
	pretty bool
}

// Process is the main function that processes the arguments and prints the output.
func Process(proc *args, writer io.Writer) error {
	if requiresHelp(proc) {
		_ = printHelpText(writer)
	}

	if err := validateArgs(proc); err != nil {
		return err
	}

	if proc.file != "" {
		proc.dirs = append(proc.dirs, proc.file)
	}

	files, err := loadFiles(proc.dirs)
	if err != nil {
		return err
	}

	tests, err := listTests(files)
	if err != nil {
		return err
	}

	marshal, err := json.Marshal(tests)
	if err != nil {
		return fmt.Errorf("%s: %w", errUnknown, err)
	}

	if proc.pretty {
		if err := prettyPrint(marshal, writer); err != nil {
			return err
		}
	} else {
		_, _ = writer.Write(marshal)
	}

	return nil
}

// validateArgs validates the arguments provided by the user.
func validateArgs(args *args) error {
	if args.file != "" && len(args.dirs) > 0 {
		return errPathIssue
	}

	if args.file != "" {
		stat, err := os.Stat(args.file)
		if err != nil {
			return fmt.Errorf("%s: %w", errUnknown, err)
		}

		if stat.IsDir() {
			return errNotAFile
		}
	}

	return nil
}

// requiresHelp checks if the user has requested help and not provided any required arguments.
func requiresHelp(proc *args) bool {
	return (len(proc.dirs) == 0 && proc.file == "") || proc.help
}

// loadFiles loads all the go files in the given paths.
func loadFiles(dirs []string) ([]string, error) {
	var testFiles []string

	for _, dir := range dirs {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() && filepath.Ext(path) == ".go" && strings.HasSuffix(path, "_test.go") {
				testFiles = append(testFiles, path)
			}

			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errUnknown, err)
		}
	}

	return testFiles, nil
}

// listTests lists all the tests in the given go test files.
func listTests(files []string) ([]Detail, error) {
	var tests []Detail

	for _, testFile := range files {
		fileAbsPath, err := filepath.Abs(testFile)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errUnknown, err)
		}

		fileName := filepath.Base(testFile)

		set := token.NewFileSet()

		parseFile, err := parser.ParseFile(set, testFile, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errUnknown, err)
		}

		for _, obj := range parseFile.Scope.Objects {
			if obj.Kind == ast.Fun {
				if strings.HasPrefix(obj.Name, "Test") ||
					strings.HasPrefix(obj.Name, "Example") ||
					strings.HasPrefix(obj.Name, "Benchmark") {
					tests = append(tests, Detail{
						Name:         obj.Name,
						Pos:          obj.Pos(),
						Line:         set.Position(obj.Pos()).Line,
						FileName:     fileName,
						RelativePath: testFile,
						AbsolutePath: fileAbsPath,
					})
				}
			}
		}
	}

	// sort the tests by name
	sort.Slice(tests, func(i, j int) bool {
		return strings.Compare(tests[i].Name, tests[j].Name) < 0
	})

	return tests, nil
}

// prettyPrint prints the given json in a pretty format.
func prettyPrint(data []byte, writer io.Writer) error {
	var prettyJSON bytes.Buffer

	err := json.Indent(&prettyJSON, data, "", "\t")
	if err != nil {
		return fmt.Errorf("%s: %w", errUnknown, err)
	}

	_, err = writer.Write(prettyJSON.Bytes())

	return err
}

// printHelpText prints the help text for the program.
func printHelpText(writer io.Writer) error {
	_, err := writer.Write([]byte(`gotest-ls provides a list of all tests in a package or a file in JSON format.

Usage:
  gotest-ls [flags] [directories]

Examples:
	gotest-ls .
 	gotest-ls -p ./cmd
 	gotest-ls -p ./cmd ./pkg
 	gotest-ls -f ./pkg/random_test.go
 	gotest-ls -p -f ./pkg/random_test.go

Flags:
  -f, --file string   Path to a file, cannot be used with directories
  -h, --help          help for gotest-ls
  -p, --pretty        Pretty print the output in JSON format`))

	return err
}
