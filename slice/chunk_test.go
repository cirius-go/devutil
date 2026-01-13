package slice_test

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cirius-go/devutil/slice"
)

func TestChunk(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		chunkSize int
		want      [][]int
	}{
		{
			name:      "nil input",
			input:     nil,
			chunkSize: 2,
			want:      nil,
		},
		{
			name:      "empty input",
			input:     []int{},
			chunkSize: 2,
			want:      [][]int{},
		},
		{
			name:      "exact split",
			input:     []int{1, 2, 3, 4},
			chunkSize: 2,
			want:      [][]int{{1, 2}, {3, 4}},
		},
		{
			name:      "uneven split",
			input:     []int{1, 2, 3, 4, 5},
			chunkSize: 2,
			want:      [][]int{{1, 2}, {3, 4}, {5}},
		},
		{
			name:      "single chunk",
			input:     []int{1, 2, 3},
			chunkSize: 5,
			want:      [][]int{{1, 2, 3}},
		},
		{
			name:      "chunk size 1",
			input:     []int{1, 2, 3},
			chunkSize: 1,
			want:      [][]int{{1}, {2}, {3}},
		},
		{
			name:      "invalid chunk size",
			input:     []int{1, 2},
			chunkSize: 0,
			want:      [][]int{{1}, {2}}, // defaults to 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slice.Chunk(tt.input, tt.chunkSize)
			if !slicesEqual2D(got, tt.want) {
				t.Errorf("Chunk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func slicesEqual2D[T comparable](a, b [][]T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func TestForEachChunk(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Test case 1: Sequential processing
	t.Run("sequential", func(t *testing.T) {
		var processed []int
		err := slice.ForEachChunk(input, 3, 1, func(chunk []int) error {
			processed = append(processed, chunk...)
			return nil
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(processed) != 10 {
			t.Errorf("expected 10 processed items, got %d", len(processed))
		}
		// In sequential mode, order should be preserved exactly as appended in blocks
		if !slicesEqual(processed, input) {
			t.Errorf("processed items mismatch order or content")
		}
	})

	// Test case 2: Concurrent processing
	t.Run("concurrent", func(t *testing.T) {
		var mu sync.Mutex
		var processed []int
		err := slice.ForEachChunk(input, 2, 3, func(chunk []int) error {
			time.Sleep(10 * time.Millisecond) // Simulate work
			mu.Lock()
			processed = append(processed, chunk...)
			mu.Unlock()
			return nil
		})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(processed) != 10 {
			t.Errorf("expected 10 processed items, got %d", len(processed))
		}
		// Order isn't guaranteed, so sort before comparing
		sort.Ints(processed)
		if !slicesEqual(processed, input) {
			t.Errorf("processed items content mismatch")
		}
	})

	// Test case 3: Error handling
	t.Run("error handling", func(t *testing.T) {
		errExpected := errors.New("oops")
		err := slice.ForEachChunk(input, 2, 2, func(chunk []int) error {
			if chunk[0] == 5 { // 3rd chunk {5, 6}
				return errExpected
			}
			return nil
		})
		if err != errExpected {
			t.Errorf("expected error %v, got %v", errExpected, err)
		}
	})

	// Test case 4: Concurrency control
	t.Run("concurrency limit", func(t *testing.T) {
		var active int32
		maxActive := int32(0)
		concurrency := 2
		// large input to ensure we fill the semaphore
		largeInput := make([]int, 20)

		err := slice.ForEachChunk(largeInput, 1, concurrency, func(chunk []int) error {
			current := atomic.AddInt32(&active, 1)
			if current > maxActive {
				// Use CAS or just careful observation. Since strict correctness of this max check
				// in test might be flaky if not locked, but here we just want to ensure we don't EXCEED.
				// Actually, atomic load/store for maxActive isn't enough to prevent race on the 'if' check
				// but for "Did we exceed?" we can just check if current > concurrency.
				if int(current) > concurrency {
					return fmt.Errorf("concurrency exceeded limit: %d > %d", current, concurrency)
				}
			}
			time.Sleep(10 * time.Millisecond)
			atomic.AddInt32(&active, -1)
			return nil
		})
		if err != nil {
			t.Errorf("concurrency check failed: %v", err)
		}
	})
}

func slicesEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
