package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListAllTestsForGivenFile(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer

	err := Process(&args{
		dirs: []string{"./tests"},
	}, &buffer)
	require.NoError(t, err)

	// This is to make it work on any OS and any folder
	pwd, err := os.Getwd()
	require.NoError(t, err)

	fmt.Println(buffer.String())

	require.JSONEq(t,
		strings.ReplaceAll(`[{"name":"BenchmarkSomething","fileName":"benchmark_test.go","relativePath":"tests/benchmark_test.go","absolutePath":"##PATH##/tests/benchmark_test.go","line":5,"pos":44},{"name":"Example_something","fileName":"example_test.go","relativePath":"tests/example_test.go","absolutePath":"##PATH##/tests/example_test.go","line":5,"pos":40},{"name":"Test/5_+_5_=_10","fileName":"table_test.go","relativePath":"tests/table_test.go","absolutePath":"##PATH##/tests/table_test.go","line":23,"pos":265},{"name":"Test/5_-_5_=_0","fileName":"table_test.go","relativePath":"tests/table_test.go","absolutePath":"##PATH##/tests/table_test.go","line":30,"pos":355},{"name":"Test/mixed_subtest_1","fileName":"table_test.go","relativePath":"tests/table_test.go","absolutePath":"##PATH##/tests/table_test.go","line":12,"pos":111},{"name":"Test/mixed_test_2","fileName":"table_test.go","relativePath":"tests/table_test.go","absolutePath":"##PATH##/tests/table_test.go","line":48,"pos":635},{"name":"TestSomething","fileName":"sample_test.go","relativePath":"tests/sample_test.go","absolutePath":"##PATH##/tests/sample_test.go","line":7,"pos":49},{"name":"Test_subTestPattern/subtest","fileName":"subtest_test.go","relativePath":"tests/subtest_test.go","absolutePath":"##PATH##/tests/subtest_test.go","line":10,"pos":121},{"name":"Test_subTestPattern/subtest_2","fileName":"subtest_test.go","relativePath":"tests/subtest_test.go","absolutePath":"##PATH##/tests/subtest_test.go","line":15,"pos":193}]`,
			"##PATH##", pwd),
		buffer.String())
}

func Example_errorIfFileAndDirectoryBothAreProvided() {
	cmd := exec.Command("go", "run", "main.go", "-p", "-f", "./tests/sample_test.go", "./tests")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err == nil {
		panic("expected error")
	}

	// Output: ERROR: cannot specify both a file and a directory
}

func Example_errorIfFileProvidedIsDirectory() {
	cmd := exec.Command("go", "run", "main.go", "-p", "-f", "./tests")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err == nil {
		panic("expected error")
	}

	// Output: ERROR: required file, provided directory
}

func Test_process(t *testing.T) {
	t.Parallel()

	pwd, err := os.Getwd()
	require.NoError(t, err)

	tests := []struct {
		name        string
		args        args
		wantErr     bool
		errExpected string
		checks      func(t *testing.T, got string)
	}{
		{
			name: "should return error if file and directory both are provided",
			args: args{
				pretty: false,
				file:   "./tests/sample_test.go",
				dirs:   []string{"./tests"},
			},
			wantErr:     true,
			errExpected: errPathIssue.Error(),
		},
		{
			name: "should return error if file provided is directory",
			args: args{
				pretty: false,
				file:   "./tests",
			},
			wantErr:     true,
			errExpected: errNotAFile.Error(),
		},
		{
			name: "should return the test details in a file",
			args: args{
				pretty: false,
				file:   "./tests/sample_test.go",
			},
			checks: func(t *testing.T, got string) {
				t.Helper()

				require.JSONEq(t, fmt.Sprintf(`[{"name":"TestSomething","fileName":"sample_test.go","relativePath":"sample_test.go","absolutePath":"%s/tests/sample_test.go","line":7,"pos":49}]`, pwd),
					got)
			},
		},
		{
			name: "should return the test details in a file with pretty flag",
			args: args{
				pretty: true,
				file:   "./tests/sample_test.go",
			},
			checks: func(t *testing.T, got string) {
				t.Helper()

				require.JSONEq(t, fmt.Sprintf(`[
	{
		"name": "TestSomething",
		"fileName": "sample_test.go",
		"relativePath": "sample_test.go",
		"absolutePath": "%s/tests/sample_test.go",
		"line": 7,
		"pos": 49
	}
]`, pwd), got)
			},
		},
		{
			name: "should also return subtests and table tests",
			args: args{
				pretty: false,
				file:   "./tests/table_test.go",
			},
			checks: func(t *testing.T, got string) {
				t.Helper()

				require.JSONEq(t, strings.ReplaceAll(`[{"name":"Test/5_+_5_=_10","fileName":"table_test.go","relativePath":"table_test.go","absolutePath":"##PATH##/tests/table_test.go","line":23,"pos":265},{"name":"Test/5_-_5_=_0","fileName":"table_test.go","relativePath":"table_test.go","absolutePath":"##PATH##/tests/table_test.go","line":30,"pos":355},{"name":"Test/mixed_subtest_1","fileName":"table_test.go","relativePath":"table_test.go","absolutePath":"##PATH##/tests/table_test.go","line":12,"pos":111},{"name":"Test/mixed_test_2","fileName":"table_test.go","relativePath":"table_test.go","absolutePath":"##PATH##/tests/table_test.go","line":48,"pos":635}]`, "##PATH##", pwd), got)
			},
		},
		{
			name: "should show help if no arguments are provided",
			args: args{},
			checks: func(t *testing.T, got string) {
				t.Helper()

				require.Equal(t, `gotest-ls provides a list of all tests in a package or a file in JSON format.

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
  -p, --pretty        Pretty print the output in JSON format
`, got)
			},
		},
		{
			name: "return error if directory does not exist",
			args: args{
				dirs: []string{"./false-directory"},
			},
			wantErr:     true,
			errExpected: errUnknown.Error() + ": lstat ./false-directory: no such file or directory",
		},
		{
			name: "return error if there is no test in the directory",
			args: args{
				dirs: []string{"./dead-tests"},
			},
			checks: func(t *testing.T, got string) {
				t.Helper()

				require.Equal(t, "No tests found\n", got)
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			writer := &bytes.Buffer{}

			err := Process(&tt.args, writer)

			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, tt.errExpected, err.Error())
			} else {
				require.NoError(t, err)
			}

			if tt.checks != nil {
				tt.checks(t, writer.String())
			}
		})
	}
}
