
## Go List Tests
![master](https://github.com/ninadingole/gotest-ls/actions/workflows/base.yml/badge.svg?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/ninadingole/gotest-ls.svg)](https://pkg.go.dev/github.com/ninadingole/gotest-ls)
[![codecov](https://codecov.io/gh/ninadingole/gotest-ls/branch/main/graph/badge.svg?token=9ZYKWNF6JI)](https://codecov.io/gh/ninadingole/gotest-ls)

`gotest-ls` is a tool to list tests in a Go project. It provides list of all the Tests
(`Test*`, `Benchmark*`, `Example*`) in a Go project or a go file.

The tool provides output in JSON format. The output can be used to generate a report or for other tools for analysis.

### Requirements
- `Go 1.18` or above

### Installation

```bash
go install github.com/ninadingole/gotest-ls
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
		"fileName": "benchmark_test.go",
		"relativePath": "tests/benchmark_test.go",
		"absolutePath": "/www/gotest-ls/tests/benchmark_test.go",
		"line": 5,
		"pos": 44
	},
	{
		"name": "Example_errorIfFileAndDirectoryBothAreProvided",
		"fileName": "main_test.go",
		"relativePath": "main_test.go",
		"absolutePath": "/www/gotest-ls/main_test.go",
		"line": 36,
		"pos": 2006
	},
	{
		"name": "Example_errorIfFileProvidedIsDirectory",
		"fileName": "main_test.go",
		"relativePath": "main_test.go",
		"absolutePath": "/www/gotest-ls/main_test.go",
		"line": 48,
		"pos": 2335
	},
	{
		"name": "Example_something",
		"fileName": "example_test.go",
		"relativePath": "tests/example_test.go",
		"absolutePath": "/www/gotest-ls/tests/example_test.go",
		"line": 5,
		"pos": 40
	},
	{
		"name": "Test/5_+_5_=_10",
		"fileName": "table_test.go",
		"relativePath": "tests/table_test.go",
		"absolutePath": "/www/gotest-ls/tests/table_test.go",
		"line": 23,
		"pos": 265
	},
	{
		"name": "Test/5_-_5_=_0",
		"fileName": "table_test.go",
		"relativePath": "tests/table_test.go",
		"absolutePath": "/www/gotest-ls/tests/table_test.go",
		"line": 30,
		"pos": 355
	},
	{
		"name": "Test/mixed_subtest_1",
		"fileName": "table_test.go",
		"relativePath": "tests/table_test.go",
		"absolutePath": "/www/gotest-ls/tests/table_test.go",
		"line": 12,
		"pos": 111
	},
	{
		"name": "Test/mixed_test_2",
		"fileName": "table_test.go",
		"relativePath": "tests/table_test.go",
		"absolutePath": "/www/gotest-ls/tests/table_test.go",
		"line": 48,
		"pos": 635
	},
	{
		"name": "TestListAllTestsForGivenFile",
		"fileName": "main_test.go",
		"relativePath": "main_test.go",
		"absolutePath": "/www/gotest-ls/main_test.go",
		"line": 14,
		"pos": 127
	},
	{
		"name": "TestSomething",
		"fileName": "sample_test.go",
		"relativePath": "tests/sample_test.go",
		"absolutePath": "/www/gotest-ls/tests/sample_test.go",
		"line": 7,
		"pos": 49
	},
	{
		"name": "Test_List/empty",
		"fileName": "list_test.go",
		"relativePath": "pkg/list_test.go",
		"absolutePath": "/www/gotest-ls/pkg/list_test.go",
		"line": 25,
		"pos": 357
	},
	{
		"name": "Test_List/fail_for_invalid_dir",
		"fileName": "list_test.go",
		"relativePath": "pkg/list_test.go",
		"absolutePath": "/www/gotest-ls/pkg/list_test.go",
		"line": 58,
		"pos": 1192
	},
	{
		"name": "Test_List/fail_to_parse_invalid_test_file",
		"fileName": "list_test.go",
		"relativePath": "pkg/list_test.go",
		"absolutePath": "/www/gotest-ls/pkg/list_test.go",
		"line": 64,
		"pos": 1328
	},
	{
		"name": "Test_List/parse_subtests_correctly",
		"fileName": "list_test.go",
		"relativePath": "pkg/list_test.go",
		"absolutePath": "/www/gotest-ls/pkg/list_test.go",
		"line": 70,
		"pos": 1500
	},
	{
		"name": "Test_List/single_dir",
		"fileName": "list_test.go",
		"relativePath": "pkg/list_test.go",
		"absolutePath": "/www/gotest-ls/pkg/list_test.go",
		"line": 44,
		"pos": 819
	},
	{
		"name": "Test_List/single_file",
		"fileName": "list_test.go",
		"relativePath": "pkg/list_test.go",
		"absolutePath": "/www/gotest-ls/pkg/list_test.go",
		"line": 30,
		"pos": 437
	},
	{
		"name": "Test_process/return_error_if_directory_does_not_exist",
		"fileName": "main_test.go",
		"relativePath": "main_test.go",
		"absolutePath": "/www/gotest-ls/main_test.go",
		"line": 164,
		"pos": 5681
	},
	{
		"name": "Test_process/return_error_if_there_is_no_test_in_the_directory",
		"fileName": "main_test.go",
		"relativePath": "main_test.go",
		"absolutePath": "/www/gotest-ls/main_test.go",
		"line": 172,
		"pos": 5920
	},
	{
		"name": "Test_process/should_also_return_subtests_and_table_tests",
		"fileName": "main_test.go",
		"relativePath": "main_test.go",
		"absolutePath": "/www/gotest-ls/main_test.go",
		"line": 127,
		"pos": 4167
	},
	{
		"name": "Test_process/should_return_error_if_file_and_directory_both_are_provided",
		"fileName": "main_test.go",
		"relativePath": "main_test.go",
		"absolutePath": "/www/gotest-ls/main_test.go",
		"line": 74,
		"pos": 2872
	},
	{
		"name": "Test_process/should_return_error_if_file_provided_is_directory",
		"fileName": "main_test.go",
		"relativePath": "main_test.go",
		"absolutePath": "/www/gotest-ls/main_test.go",
		"line": 84,
		"pos": 3124
	},
	{
		"name": "Test_process/should_return_the_test_details_in_a_file",
		"fileName": "main_test.go",
		"relativePath": "main_test.go",
		"absolutePath": "/www/gotest-ls/main_test.go",
		"line": 93,
		"pos": 3317
	},
	{
		"name": "Test_process/should_return_the_test_details_in_a_file_with_pretty_flag",
		"fileName": "main_test.go",
		"relativePath": "main_test.go",
		"absolutePath": "/www/gotest-ls/main_test.go",
		"line": 106,
		"pos": 3722
	},
	{
		"name": "Test_process/should_show_help_if_no_arguments_are_provided",
		"fileName": "main_test.go",
		"relativePath": "main_test.go",
		"absolutePath": "/www/gotest-ls/main_test.go",
		"line": 139,
		"pos": 5054
	},
	{
		"name": "Test_subTestPattern/subtest",
		"fileName": "subtest_test.go",
		"relativePath": "tests/subtest_test.go",
		"absolutePath": "/www/gotest-ls/tests/subtest_test.go",
		"line": 10,
		"pos": 121
	},
	{
		"name": "Test_subTestPattern/subtest_2",
		"fileName": "subtest_test.go",
		"relativePath": "tests/subtest_test.go",
		"absolutePath": "/www/gotest-ls/tests/subtest_test.go",
		"line": 15,
		"pos": 193
	}
]
```
