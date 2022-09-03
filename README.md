
## Go List Tests

`gotest-ls` is a tool to list tests in a Go project. It provides list of all the Tests
(`Test*`, `Benchmark*`, `Example*`) in a Go project or a go file.

The tool provides output in JSON format. The output can be used to generate a report or for other tools for analysis.

### Requirements
- `Go 1.18` or above

### Installation

```bash
go get -u github.com/ninadingole/gotest-ls
```

### Usage

```bash
gotest-ls [flags] [directories]

gotest-ls .
gotest-ls -p ./cmd
gotest-ls -p ./cmd ./pkg
gotest-ls -f ./pkg/random_test.go
gotest-ls -p -f ./pkg/random_test.go

```

### Flags

```bash
  -h, --help                help for gotest-ls
  -f, --file    string      file to list tests from
  -p, --pretty  bool        pretty print the json output
```

### Output

```bash
$> gotest-ls -p .                        
[
        {
                "name": "BenchmarkSomething",
                "fileName": "sample_test.go",
                "relativePath": "tests/sample_test.go",
                "absolutePath": "www/gotest-ls/tests/sample_test.go",
                "line": 14,
                "pos": 156
        },
        {
                "name": "Example_errorIfFileAndDirectoryBothAreProvided",
                "fileName": "main_test.go",
                "relativePath": "main_test.go",
                "absolutePath": "www/gotest-ls/main_test.go",
                "line": 37,
                "pos": 1090
        },
        {
                "name": "Example_errorIfFileProvidedIsDirectory",
                "fileName": "main_test.go",
                "relativePath": "main_test.go",
                "absolutePath": "www/gotest-ls/main_test.go",
                "line": 49,
                "pos": 1419
        },
        {
                "name": "Example_something",
                "fileName": "sample_test.go",
                "relativePath": "tests/sample_test.go",
                "absolutePath": "www/gotest-ls/tests/sample_test.go",
                "line": 20,
                "pos": 250
        },
        {
                "name": "TestListAllTestsForGivenFile",
                "fileName": "main_test.go",
                "relativePath": "main_test.go",
                "absolutePath": "www/gotest-ls/main_test.go",
                "line": 14,
                "pos": 126
        },
        {
                "name": "TestSomething",
                "fileName": "sample_test.go",
                "relativePath": "tests/sample_test.go",
                "absolutePath": "www/gotest-ls/tests/sample_test.go",
                "line": 8,
                "pos": 56
        }
]

```
