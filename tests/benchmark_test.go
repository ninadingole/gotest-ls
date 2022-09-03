package tests_test

import "testing"

func BenchmarkSomething(b *testing.B) {
	b.Skipf("Skipping...")

	b.Log("Benchmarking...")
}
