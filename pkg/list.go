// Package pkg contains the core logic of the gotest-ls tool which finds all the go test files.
// and using ast package, it lists all the tests in the given files.
package pkg

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
)

// testType represents the type of test function.
type testType int

const (
	testTypeNone testType = iota
	testTypeSubTest
	testTypeTableTest
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

// subTestDetail returns the testname and the position of the subtest in the file.
type subTestDetail struct {
	name string
	pos  token.Pos
}

// List returns all the go test files in the given directories or a given file.
// It returns an error if the given directories are invalid.
// It returns an empty slice if no tests are found.
// The returned slice is sorted by the test name.
func List(fileOrDirs []string) ([]TestDetail, error) {
	files, err := loadFiles(fileOrDirs)
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
func loadFiles(dirs []string) (map[string][]string, error) {
	testFiles := make(map[string][]string)

	for _, dir := range dirs {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && filepath.Ext(path) == ".go" && strings.HasSuffix(path, "_test.go") {
				testFiles[dir] = append(testFiles[dir], path)
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
func listTests(files map[string][]string) ([]TestDetail, error) { //nolint: gocognit
	var tests []TestDetail

	for dir, testFiles := range files {
		for _, testFile := range testFiles {
			set := token.NewFileSet()

			parseFile, err := parser.ParseFile(set, testFile, nil, parser.AllErrors)
			if err != nil {
				return nil, err
			}

			ast.Inspect(parseFile, func(n ast.Node) bool {
				fnDecl, ok := n.(*ast.FuncDecl)
				if !ok {
					return true
				}

				if !isGolangTest(fnDecl) {
					return false
				}
				isSubTest := false
				isParallel := false

				for i, v := range fnDecl.Body.List {
					if !isParallel {
						isParallel = isParallelTest(v)
					}
					switch identifyTestType(v) {
					case testTypeSubTest:
						isSubTest = true

						if test := findSubTestName(v); test != nil {
							tests = append(tests, buildTestDetail(fnDecl.Name.String(), test.name, dir, testFile, set, test.pos))
						}

					case testTypeTableTest:
						isSubTest = true
						testNameFieldInStruct := findTableTestNameField(v)

						if testNameFieldInStruct != "" {
							for j := i; j > 0; j-- {
								if ttDetails := parseTableTestStructsIfAny(fnDecl.Body.List[j], testNameFieldInStruct); ttDetails != nil {
									for _, ttDetail := range ttDetails {
										tests = append(tests, buildTestDetail(fnDecl.Name.String(), ttDetail.name, dir, testFile, set, ttDetail.pos))
									}
								}
							}
						}
					case testTypeNone:
						continue
					}
				}
				if !isSubTest {
					tests = append(tests, buildTestDetail(fnDecl.Name.String(), "", dir, testFile, set, fnDecl.Name.Pos()))
				}

				return true
			})
		}
	}

	// sort the tests by name
	sort.Slice(tests, func(i, j int) bool {
		return strings.Compare(tests[i].Name, tests[j].Name) < 0
	})

	return tests, nil
}

func isParallelTest(v ast.Stmt) bool {
	exprStmt, ok := v.(*ast.ExprStmt)
	if !ok {
		return false
	}

	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return false
	}

	if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
		if selectorExpr.Sel.Name == "Parallel" {
			return true
		}
	}

	return false
}

// isGolangTest checks if the function name starts with golang test standards
// it checks for `Test`, `Example`, `Benchmark` or `Fuzz` prefixes in a function name.
// Other than test functions all the other functions are ignored.
func isGolangTest(obj *ast.FuncDecl) bool {
	name := obj.Name.String()
	testPrefixes := []string{"Test", "Example", "Benchmark", "Fuzz"}

	for _, prefix := range testPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}

	return false
}

// identifyTestType identifies the type of the test based on the given ast node.
// it looks for `t.Run` function in the test function body. If the test contains subtests then it returns
// testTypeSubTest. If the test contains table tests then it returns testTypeTableTest.
// Otherwise, it returns testTypeNone.
func identifyTestType(v ast.Stmt) testType {
	switch expr := v.(type) {
	case *ast.ExprStmt:
		if callExpr, ok := expr.X.(*ast.CallExpr); ok {
			if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
				if selectorExpr.Sel.Name == "Run" {
					return testTypeSubTest
				}
			}
		}
	case *ast.RangeStmt:
		for _, v := range expr.Body.List {
			if typ := identifyTestType(v); typ == testTypeSubTest {
				return testTypeTableTest
			}
		}
	}

	return testTypeNone
}

// findSubTestName finds the name of the subtest in the given ast node.
// it looks for `t.Run` function in the test function body. If the test contains subtests then it returns
// the name of the subtest.
// A test would look like this in the source code.
//
//	func Test_subTestPattern(t *testing.T) {
//		t.Parallel()
//
//		msg := "Hello, world!"
//
//		t.Run("subtest", func(t *testing.T) {
//			t.Parallel()
//			t.Log(msg)
//		})
//
//		t.Run("subtest 2", func(t *testing.T) {
//			t.Parallel()
//			t.Log("This is a subtest")
//		})
//	}
func findSubTestName(v ast.Stmt) *subTestDetail {

	expr, ok := v.(*ast.ExprStmt)
	if !ok {
		return nil
	}

	callExpr, ok := expr.X.(*ast.CallExpr)
	if !ok {
		return nil
	}

	if basic, ok := callExpr.Args[0].(*ast.BasicLit); ok {
		return &subTestDetail{
			name: basic.Value,
			pos:  callExpr.Pos(),
		}
	}

	return nil
}

// buildTestDetail returns the TestDetail object with the information received from the given parameters.
func buildTestDetail(parent string, subTest string, dir string, file string, set *token.FileSet, pos token.Pos) TestDetail {
	fileAbsPath, err := filepath.Abs(file)
	if err != nil {
		panic(fmt.Errorf("failed to get absolute path of file %s: %w", file, err))
	}

	fileName := filepath.Base(file)

	relativePath, err := filepath.Rel(filepath.Dir(dir), file)
	if err != nil {
		panic(fmt.Errorf("failed to get relative path of file %s: %w", file, err))
	}

	detail := TestDetail{
		Name:         parent,
		FileName:     fileName,
		RelativePath: relativePath,
		AbsolutePath: fileAbsPath,
		Line:         set.Position(pos).Line,
		Pos:          pos,
	}

	if subTest != "" {
		detail.Name = fmt.Sprintf("%s/%s", parent,
			strings.ReplaceAll(strings.ReplaceAll(subTest, "\"", ""), " ", "_"))
	}

	return detail
}

// findTableTestNameField returns the name of the field in the table test struct which contains the test name.
// it looks for the field used in `t.Run` inside the for-loop of a table test and returns the name of the parameter
// from the struct that is used to populate the test name.
// A typical table test range function would look like this in the source code.
//
//	for _, tt := range tests {
//			tt := tt
//			t.Run(tt.name, func(t *testing.T) {
//				t.Parallel()
//
//				if got := tt.calc(); got != tt.want {
//					t.Errorf("got %d, want %d", got, tt.want)
//				}
//			})
//		}
func findTableTestNameField(v ast.Stmt) string {
	rangeStmt, ok := v.(*ast.RangeStmt)
	if !ok {
		return ""
	}

	if len(rangeStmt.Body.List) == 0 {
		return ""
	}

	for _, stmt := range rangeStmt.Body.List {
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}
		callExpr, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			continue
		}
		selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		ident, ok := selectorExpr.X.(*ast.Ident)
		if !ok {
			continue
		}
		if ident.Name == "t" && selectorExpr.Sel.Name == "Run" {
			if sExpr, ok := callExpr.Args[0].(*ast.SelectorExpr); ok {
				return strings.ReplaceAll(sExpr.Sel.Name, "\"", "")
			}
		}

	}

	return ""
}

// parseTableTestStructsIfAny parses the struct array in the table test and returns the value of the field that
// will be passed to `t.Run` function when the test is run.
func parseTableTestStructsIfAny(v ast.Stmt, fieldName string) []subTestDetail {
	var values []subTestDetail

	assignStmt, ok := v.(*ast.AssignStmt)
	if !ok {
		return nil
	}
	for _, expr := range assignStmt.Rhs {
		cmpsLit, ok := expr.(*ast.CompositeLit)
		if !ok {
			continue
		}
		for _, elt := range cmpsLit.Elts {
			compositeLit, ok := elt.(*ast.CompositeLit)
			if !ok {
				continue
			}
			for _, elt := range compositeLit.Elts {
				kvExpr, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				key, ok := kvExpr.Key.(*ast.Ident)
				if !ok {
					continue
				}

				if key.Name == fieldName {
					if value, ok := kvExpr.Value.(*ast.BasicLit); ok {
						values = append(values,
							subTestDetail{
								name: value.Value,
								pos:  key.Pos(),
							})
					}
				}

			}
		}
	}

	return values
}
