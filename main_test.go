package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/ninadingole/gotest-ls/pkg"
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

	require.JSONEq(t,
		strings.ReplaceAll(`[{"name":"BenchmarkSomething","fileName":"benchmark_test.go","relativePath":"tests/benchmark_test.go","absolutePath":"##PATH##/tests/benchmark_test.go","line":5,"pos":44},{"name":"Example_something","fileName":"example_test.go","relativePath":"tests/example_test.go","absolutePath":"##PATH##/tests/example_test.go","line":5,"pos":40},{"name":"TestSomething","fileName":"sample_test.go","relativePath":"tests/sample_test.go","absolutePath":"##PATH##/tests/sample_test.go","line":7,"pos":49}]`,
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
		expected    string
		errExpected string
		checks      func(t *testing.T, got []pkg.TestDetail)
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
			expected: fmt.Sprintf(`[{"name":"TestSomething","fileName":"sample_test.go","relativePath":"./tests/sample_test.go","absolutePath":"%s/tests/sample_test.go","line":7,"pos":49}]`, pwd),
		},
		{
			name: "should return the test details in a file with pretty flag",
			args: args{
				pretty: true,
				file:   "./tests/sample_test.go",
			},
			expected: fmt.Sprintf(`[
	{
		"name": "TestSomething",
		"fileName": "sample_test.go",
		"relativePath": "./tests/sample_test.go",
		"absolutePath": "%s/tests/sample_test.go",
		"line": 7,
		"pos": 49
	}
]`, pwd),
		},
		{
			name: "should show help if no arguments are provided",
			args: args{},
			expected: `gotest-ls provides a list of all tests in a package or a file in JSON format.

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
  -p, --pretty        Pretty print the output in JSON format`,
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
			expected: "No tests found",
			wantErr:  false,
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

			fmt.Println(writer.String())
			require.Equal(t, tt.expected, writer.String())
		})
	}
}
