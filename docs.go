// gotest-ls provides a list of all tests in a package or a file in JSON format.
//
// Usage:
//
//	gotest-ls [flags] [directories...]
//
// Examples:
//
//	gotest-ls .
//	gotest-ls -p ./cmd
//	gotest-ls -p ./cmd ./pkg
//	gotest-ls -f ./pkg/random_test.go
//	gotest-ls -p -f ./pkg/random_test.go
//
// Flags:
//
//	-f, --file string   Path to a file, cannot be used with directories
//	-h, --help          help for gotest-ls
//	-p, --pretty        Pretty print the output in JSON format
package main
