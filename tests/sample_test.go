package tests_test

import (
	"fmt"
	"testing"
)

func TestSomething(t *testing.T) {
	t.Parallel()
	t.Skipf("Skipping...")
	t.Log("Hello, world!")
}

func BenchmarkSomething(b *testing.B) {
	b.Skipf("Skipping...")

	b.Log("Benchmarking...")
}

func Example_something() {
	fmt.Println("Example!")
	// Output: Example!
}
