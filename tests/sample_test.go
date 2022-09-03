package tests_test

import (
	"testing"
)

func TestSomething(t *testing.T) {
	t.Parallel()
	t.Skipf("Skipping...")
	t.Log("Hello, world!")
}
