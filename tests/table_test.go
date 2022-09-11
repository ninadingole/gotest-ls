package tests_test

import (
	"testing"
)

func Test(t *testing.T) {
	t.Parallel()

	msg := "Hello, world!"

	t.Run("mixed subtest 1", func(t *testing.T) {
		t.Parallel()
		t.Log(msg)
	})

	tests := []struct {
		name string
		calc func() int
		want int
	}{
		{
			name: "5 + 5 = 10",
			calc: func() int {
				return 5 + 5
			},
			want: 10,
		},
		{
			name: "5 - 5 = 0",
			calc: func() int {
				return 5 - 2
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.calc(); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}

	t.Run("mixed test 2", func(t *testing.T) {
		t.Parallel()
		t.Log("This is a subtest")
	})
}
