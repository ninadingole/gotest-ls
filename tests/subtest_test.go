package tests_test

import "testing"

func Test_subTestPattern(t *testing.T) {
	t.Parallel()

	msg := "Hello, world!"

	t.Run("subtest", func(t *testing.T) {
		t.Parallel()
		t.Log(msg)
	})

	t.Run("subtest 2", func(t *testing.T) {
		t.Parallel()
		t.Log("This is a subtest")
	})
}
