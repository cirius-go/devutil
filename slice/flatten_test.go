package slice_test

import (
	"testing"

	"github.com/cirius-go/devutil/slice"
)

func TestFlatten(t *testing.T) {
	tests := []struct {
		name  string
		input [][]int
		want  []int
	}{
		{
			name:  "nil input",
			input: nil,
			want:  nil,
		},
		{
			name:  "empty input",
			input: [][]int{},
			want:  []int{},
		},
		{
			name:  "slice of empty slices",
			input: [][]int{{}, {}},
			want:  []int{},
		},
		{
			name:  "normal case",
			input: [][]int{{1}, {2, 3}, {4}},
			want:  []int{1, 2, 3, 4},
		},
		{
			name:  "mixed empty and non-empty",
			input: [][]int{{1}, {}, {2}},
			want:  []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slice.Flatten(tt.input)
			if !slicesEqual(got, tt.want) {
				t.Errorf("Flatten() = %v, want %v", got, tt.want)
			}
		})
	}
}
