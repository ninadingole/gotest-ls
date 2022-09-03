// Package pkg contains the core logic of the gotest-ls tool which finds all the go test files.
// and using ast package, it lists all the tests in the given files.
package pkg

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

// TestDetail is a struct that contains the details of a single test.
// It contains the name of the test, the line number, the file name, the relative path and the absolute path.
// It also contains the token position (token.Pos) of the test in the file.
type TestDetail struct {
	Name         string    `json:"name"`
	FileName     string    `json:"fileName"`
	RelativePath string    `json:"relativePath"`
	AbsolutePath string    `json:"absolutePath"`
	Line         int       `json:"line"`
	Pos          token.Pos `json:"pos"`
}

// List returns all the go test files in the given directories.
// It returns an error if the given directories are invalid.
// It returns an empty slice if no tests are found.
// The returned slice is sorted by the test name.
func List(dirs []string) ([]TestDetail, error) {
	files, err := loadFiles(dirs)
	if err != nil {
		return nil, err
	}

	tests, err := listTests(files)
	if err != nil {
		return nil, err
	}

	return tests, nil
}

// loadFiles loads all the go files in the given paths.
func loadFiles(dirs []string) ([]string, error) {
	var testFiles []string

	for _, dir := range dirs {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && filepath.Ext(path) == ".go" && strings.HasSuffix(path, "_test.go") {
				testFiles = append(testFiles, path)
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return testFiles, nil
}

// listTests lists all the tests in the given go test files.
func listTests(files []string) ([]TestDetail, error) {
	var tests []TestDetail

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
					tests = append(tests, TestDetail{
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
