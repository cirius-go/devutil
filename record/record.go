package record

import (
	"cmp"
	"slices"
)

// Keys returns a slice of keys from the map.
// The order of keys is not guaranteed.
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of values from the map.
// The order of values matches the order of keys (which is random).
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Clone creates a shallow copy of the map.
func Clone[K comparable, V any](m map[K]V) map[K]V {
	if m == nil {
		return nil
	}
	clone := make(map[K]V, len(m))
	for k, v := range m {
		clone[k] = v
	}
	return clone
}

// Merge merges multiple maps into a new map.
// Keys from later maps override keys from earlier maps.
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	size := 0
	for _, m := range maps {
		size += len(m)
	}
	result := make(map[K]V, size)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// Filter returns a new map containing only the entries that satisfy the predicate.
func Filter[K comparable, V any](m map[K]V, predicate func(K, V) bool) map[K]V {
	if m == nil {
		return nil
	}
	result := make(map[K]V)
	for k, v := range m {
		if predicate(k, v) {
			result[k] = v
		}
	}
	return result
}

// MapValues transforms the values of a map using a mapper function.
func MapValues[K comparable, InV, OutV any](m map[K]InV, mapper func(InV) OutV) map[K]OutV {
	if m == nil {
		return nil
	}
	result := make(map[K]OutV, len(m))
	for k, v := range m {
		result[k] = mapper(v)
	}
	return result
}

// ToSet creates a map where the keys are the elements of the slice and values are struct{}{}.
func ToSet[K comparable](input []K) map[K]struct{} {
	if input == nil {
		return nil
	}
	result := make(map[K]struct{}, len(input))
	for _, item := range input {
		result[item] = struct{}{}
	}
	return result
}

// Associate creates a map from a slice using a transform function that returns a key-value pair.
func Associate[T any, K comparable, V any](input []T, transform func(T) (K, V)) map[K]V {
	if input == nil {
		return nil
	}
	result := make(map[K]V, len(input))
	for _, item := range input {
		k, v := transform(item)
		result[k] = v
	}
	return result
}

// SortedKeys returns a slice of the map's keys, sorted.
func SortedKeys[K cmp.Ordered, V any](m map[K]V) []K {
	keys := Keys(m)
	slices.Sort(keys)
	return keys
}

// SortedValues returns a slice of the map's values, sorted.
func SortedValues[K comparable, V cmp.Ordered](m map[K]V) []V {
	values := Values(m)
	slices.Sort(values)
	return values
}
