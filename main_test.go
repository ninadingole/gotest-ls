package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListAllTestsForGivenFile(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer

	cmd := exec.Command("go", "run", "main.go", "-f", "./tests/sample_test.go")
	cmd.Stdout = &buffer
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	assert.NoError(t, err)

	fmt.Println(buffer.String())

	// This is to make it work on any OS and any folder
	pwd, err := os.Getwd()
	assert.NoError(t, err)

	assert.JSONEq(t,
		strings.ReplaceAll(`[{"name":"BenchmarkSomething","fileName":"sample_test.go","relativePath":"./tests/sample_test.go","absolutePath":"##PATH##/tests/sample_test.go","line":14,"pos":156},{"name":"Example_something","fileName":"sample_test.go","relativePath":"./tests/sample_test.go","absolutePath":"##PATH##/tests/sample_test.go","line":20,"pos":250},{"name":"TestSomething","fileName":"sample_test.go","relativePath":"./tests/sample_test.go","absolutePath":"##PATH##/tests/sample_test.go","line":8,"pos":56}]`,
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
