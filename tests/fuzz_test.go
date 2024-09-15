package tests_test

import "testing"

func Fuzz_Sample(f *testing.F) {
	f.Fuzz(func(_ *testing.T, data []byte) {
		// Dummy fuzzing logic
		if len(data) == 0 {
			return
		}

		_ = data[0]
	})
}
