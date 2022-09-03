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
		errExpected error
		checks      func(t *testing.T, got []Detail)
	}{
		{
			name: "should return error if file and directory both are provided",
			args: args{
				pretty: false,
				file:   "./tests/sample_test.go",
				dirs:   []string{"./tests"},
			},
			wantErr:     true,
			errExpected: errPathIssue,
		},
		{
			name: "should return error if file provided is directory",
			args: args{
				pretty: false,
				file:   "./tests",
			},
			wantErr:     true,
			errExpected: errNotAFile,
		},
		{
			name: "should return the test details in a file",
			args: args{
				pretty: false,
				file:   "./tests/sample_test.go",
			},
			expected: fmt.Sprintf(`[{"name":"TestSomething","fileName":"sample_test.go","relativePath":"./tests/sample_test.go","absolutePath":"%s/tests/sample_test.go","line":7,"pos":49}]`, pwd),
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
				require.Equal(t, tt.errExpected, err)
			} else {
				require.NoError(t, err)
			}

			fmt.Println(writer.String())
			require.Equal(t, tt.expected, writer.String())
		})
	}
}
