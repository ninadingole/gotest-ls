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

	if requiresHelp() {
		printHelpText()
		os.Exit(0)
	}

	tests, err := listTests(loadFiles(*file, flag.Args()))
	if err != nil {
		panic(err)
	}

	marshal, err := json.Marshal(tests)
	if err != nil {
		panic(err)
	}

	if *pretty {
		prettyPrint(marshal)
	} else {
		fmt.Println(string(marshal))
	}
}

// requiresHelp checks if the user has requested help and not provided any required arguments.
func requiresHelp() bool {
	return (len(flag.Args()) == 0 && *file == "") || *help
}

// loadFiles loads all the go files in the given paths.
func loadFiles(file string, args []string) []string {
	if file != "" && len(args) > 0 {
		fmt.Println(errPathIssue)
		os.Exit(1)
	}

	if file != "" {
		stat, err := os.Stat(file)
		if err != nil {
			panic(err)
		}

		if stat.IsDir() {
			fmt.Println(errNotAFile)
			os.Exit(1)
		}

		return []string{file}
	}

	var testFiles []string

	for _, dir := range args {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() && filepath.Ext(path) == ".go" && strings.HasSuffix(path, "_test.go") {
				testFiles = append(testFiles, path)
			}

			return nil
		})
		if err != nil {
			panic(err)
		}
	}

	return testFiles
}

// listTests lists all the tests in the given go test files.
func listTests(files []string) ([]Detail, error) {
	var tests []Detail

	for _, testFile := range files {
		fileAbsPath, err := filepath.Abs(testFile)
		if err != nil {
			return nil, err
		}

		fileName := filepath.Base(testFile)

		set := token.NewFileSet()

		parseFile, err := parser.ParseFile(set, testFile, nil, parser.ParseComments)
		if err != nil {
			return nil, err
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
func prettyPrint(data []byte) {
	var prettyJSON bytes.Buffer

	err := json.Indent(&prettyJSON, data, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(prettyJSON.String())
}

// printHelpText prints the help text for the program.
func printHelpText() {
	fmt.Println(`gotest-ls provides a list of all tests in a package or a file in JSON format.

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
  -p, --pretty        Pretty print the output in JSON format`)
}
