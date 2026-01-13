package slice

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestCollect_ErrorSwallowing(t *testing.T) {
	input := []int{1, 2, 3}
	expectedErr := errors.New("stop error")

	_, err := Collect(input, func(c CollectorContext[int, int]) {
		_, val := c.CurrentElem()
		if val == 2 {
			c.Stop(expectedErr)
			return // Implicitly covered by Stop panic, but good for clarity if it didn't panic
		}
		c.SetValue(val)
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestCollect_CurrentResult(t *testing.T) {
	input := []int{1, 2, 3}

	_, err := Collect(input, func(c CollectorContext[int, int]) {
		idx, val := c.CurrentElem()
		if idx == 1 { // Second element
			res := c.CurrentResult()
			if len(res) != 1 {
				t.Errorf("Expected CurrentResult length 1, got %d", len(res))
			}
			if len(res) > 0 && res[0] != 1 {
				t.Errorf("Expected CurrentResult[0] to be 1, got %v", res[0])
			}
		}
		c.SetValue(val)
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestCollect_ContinueByErr(t *testing.T) {
	input := []int{1, 2, 3, 4}
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")

	_, err := Collect(input, func(c CollectorContext[int, int]) {
		_, val := c.CurrentElem()
		if val == 2 {
			c.Continue(err1)
		}
		if val == 4 {
			c.Continue(err2)
		}
		c.SetValue(val)
	})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// We expect a joined error. Ideally we check if it contains both.
	// In Go 1.20+, errors.Is or string checks work.
	// Check for SliceError type and inspect elements
	var sliceErr SliceError[int]
	if !errors.As(err, &sliceErr) {
		t.Errorf("Expected error to be SliceError, got %T", err)
	} else {
		if sliceErr.At(0) == nil || sliceErr.At(0).Error() != err1.Error() {
			t.Errorf("Expected first error %v, got %v", err1, sliceErr.At(0))
		}
		// Also verify Unwrap allows standard errors.Is checks
		if !errors.Is(err, err1) {
			t.Errorf("Expected error to contain %v", err1)
		}
		if !errors.Is(err, err2) {
			t.Errorf("Expected error to contain %v", err2)
		}
	}
}

func TestCollect_NilErrors(t *testing.T) {
	input := []int{1}

	_, err := Collect(input, func(c CollectorContext[int, int]) {
		c.SetValue(0)
		c.Continue(nil, nil)
	})

	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	_, err = Collect(input, func(c CollectorContext[int, int]) {
		c.SetValue(0)
		c.Stop(nil)
	})

	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}

func TestCollect_HandlerReturnsFalse(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}

	result, err := Collect(input, func(c CollectorContext[int, int]) {
		idx, val := c.CurrentElem()
		if idx == 2 { // Stop at 3rd element (value 3)
			c.Stop()
		}
		c.SetValue(val)
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected result length 2, got %d", len(result))
	}
	if result[0] != 1 || result[1] != 2 {
		t.Errorf("Expected result [1, 2], got %v", result)
	}
}

func TestCollect_EmptyInput(t *testing.T) {
	var input []int
	res, err := Collect(input, func(c CollectorContext[int, int]) {
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("Expected empty result, got %v", res)
	}
}

func TestCollect_NilHandler(t *testing.T) {
	input := []int{1, 2}
	res, err := Collect[int, int](input, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("Expected empty result, got %v", res)
	}
}

func TestCollect_AccessSlice(t *testing.T) {
	input := []int{10, 20}
	_, err := Collect(input, func(c CollectorContext[int, int]) {
		original := c.Slice()
		if len(original) != 2 {
			t.Errorf("Expected original slice length 2, got %d", len(original))
		}
		if original[0] != 10 || original[1] != 20 {
			t.Errorf("Expected original slice [10, 20], got %v", original)
		}
		c.SetValue(0)
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestCollect_TypeTransformation(t *testing.T) {
	input := []int{1, 2, 3}
	res, err := Collect(input, func(c CollectorContext[int, string]) {
		_, val := c.CurrentElem()
		c.SetValue(fmt.Sprintf("val-%d", val))
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := []string{"val-1", "val-2", "val-3"}
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("Expected %v, got %v", expected, res)
	}
}

func TestFilter(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	res := Filter(input, func(item int) bool {
		return item%2 == 0
	})
	expected := []int{2, 4}
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("Expected %v, got %v", expected, res)
	}
}

func TestReduce(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	res := Reduce(input, func(acc int, item int) int {
		return acc + item
	}, 0)
	expected := 15
	if res != expected {
		t.Errorf("Expected %v, got %v", expected, res)
	}
}

func TestEvery(t *testing.T) {
	if !Every([]int{2, 4, 6}, func(i int) bool { return i%2 == 0 }) {
		t.Error("Expected Every to return true for all even numbers")
	}
	if Every([]int{2, 3, 6}, func(i int) bool { return i%2 == 0 }) {
		t.Error("Expected Every to return false for mixed numbers")
	}
	if !Every([]int{}, func(i int) bool { return false }) {
		t.Error("Expected Every to return true for empty slice")
	}
}

func TestSome(t *testing.T) {
	if !Some([]int{1, 3, 6}, func(i int) bool { return i%2 == 0 }) {
		t.Error("Expected Some to return true if one even number exists")
	}
	if Some([]int{1, 3, 5}, func(i int) bool { return i%2 == 0 }) {
		t.Error("Expected Some to return false if no even numbers exist")
	}
	if Some([]int{}, func(i int) bool { return true }) {
		t.Error("Expected Some to return false for empty slice")
	}
}

func TestMap(t *testing.T) {
	input := []int{1, 2, 3}
	res := Map(input, func(i int) string {
		return fmt.Sprintf("%d", i)
	})
	expected := []string{"1", "2", "3"}
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("Expected %v, got %v", expected, res)
	}
	if Map[int, int](nil, nil) != nil {
		t.Error("Expected nil result for nil input")
	}
}

func TestFind(t *testing.T) {
	input := []int{1, 2, 3, 4}
	val, found := Find(input, func(i int) bool { return i%2 == 0 })
	if !found || val != 2 {
		t.Errorf("Expected to find 2, got %v found=%v", val, found)
	}

	val, found = Find(input, func(i int) bool { return i > 10 })
	if found || val != 0 {
		t.Errorf("Expected not found, got %v found=%v", val, found)
	}
}

func TestContains(t *testing.T) {
	input := []int{1, 2, 3}
	if !Contains(input, 2) {
		t.Error("Expected Contains to return true for existing element")
	}
	if Contains(input, 4) {
		t.Error("Expected Contains to return false for non-existing element")
	}
	if Contains[int](nil, 1) {
		t.Error("Expected Contains to return false for nil input")
	}
}
