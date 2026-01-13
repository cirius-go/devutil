package record

import (
	"reflect"
	"sort"
	"testing"
)

func TestKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	keys := Keys(m)
	sort.Strings(keys)
	expected := []string{"a", "b"}
	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Expected %v, got %v", expected, keys)
	}
}

func TestValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	vals := Values(m)
	sort.Ints(vals)
	expected := []int{1, 2}
	if !reflect.DeepEqual(vals, expected) {
		t.Errorf("Expected %v, got %v", expected, vals)
	}
}

func TestClone(t *testing.T) {
	m := map[string]int{"a": 1}
	clone := Clone(m)
	if !reflect.DeepEqual(m, clone) {
		t.Errorf("Expected clone to be equal")
	}
	clone["b"] = 2
	if _, ok := m["b"]; ok {
		t.Errorf("Clone should be a deep copy of structure (shallow values)")
	}
}

func TestMerge(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 3, "c": 4}
	merged := Merge(m1, m2)
	expected := map[string]int{"a": 1, "b": 3, "c": 4}
	if !reflect.DeepEqual(merged, expected) {
		t.Errorf("Expected %v, got %v", expected, merged)
	}
}

func TestFilter(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	filtered := Filter(m, func(k string, v int) bool {
		return v%2 != 0
	})
	expected := map[string]int{"a": 1, "c": 3}
	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("Expected %v, got %v", expected, filtered)
	}
}

func TestMapValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	mapped := MapValues(m, func(v int) string {
		if v == 1 {
			return "one"
		}
		return "two"
	})
	expected := map[string]string{"a": "one", "b": "two"}
	if !reflect.DeepEqual(mapped, expected) {
		t.Errorf("Expected %v, got %v", expected, mapped)
	}
}

func TestToSet(t *testing.T) {
	input := []string{"a", "b", "a"}
	set := ToSet(input)
	if len(set) != 2 {
		t.Errorf("Expected set size 2, got %d", len(set))
	}
	if _, ok := set["a"]; !ok {
		t.Error("Expected set to contain 'a'")
	}
	if _, ok := set["b"]; !ok {
		t.Error("Expected set to contain 'b'")
	}
}

func TestAssociate(t *testing.T) {
	input := []string{"a", "b"}
	// associate: key = item, value = item + "1"
	m := Associate(input, func(item string) (string, string) {
		return item, item + "1"
	})
	expected := map[string]string{"a": "a1", "b": "b1"}
	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Expected %v, got %v", expected, m)
	}
}

func TestSortedKeys(t *testing.T) {
	m := map[string]int{"b": 2, "a": 1, "c": 3}
	keys := SortedKeys(m)
	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Expected %v, got %v", expected, keys)
	}
}

func TestSortedValues(t *testing.T) {
	m := map[string]int{"a": 2, "b": 1, "c": 3}
	vals := SortedValues(m)
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(vals, expected) {
		t.Errorf("Expected %v, got %v", expected, vals)
	}
}
