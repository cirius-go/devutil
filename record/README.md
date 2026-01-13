# Record Package

The `record` package provides utility functions for working with Go maps (associative arrays/records). It focuses on ergonomics, transformation, and common patterns like extracting keys/values or set operations.

## Features

- **Accessors**: Easily retrieve keys and values, optionally sorted.
- **Transformations**: Filter, map values, or create maps from slices.
- **Operations**: Clone, merge, and set creation.

## Usage

### Accessors

```go
m := map[string]int{"b": 2, "a": 1}

keys := record.Keys(m)         // ["b", "a"] (order random)
sorted := record.SortedKeys(m) // ["a", "b"]
vals := record.Values(m)       // [2, 1]
```

### Transformations

```go
// Filter map
odds := record.Filter(m, func(k string, v int) bool {
    return v%2 != 0
})

// Transform values
strVals := record.MapValues(m, func(v int) string {
    return fmt.Sprintf("%d", v)
})
```

### Slice to Map Conversions

```go
input := []string{"apple", "banana"}

// Create a Set
set := record.ToSet(input) 
// map[string]struct{}{"apple":{}, "banana":{}}

// Associate (Transform to Map)
lengthMap := record.Associate(input, func(s string) (string, int) {
    return s, len(s)
})
// map[string]int{"apple": 5, "banana": 6}
```

### Operations

```go
// Merge multiple maps
merged := record.Merge(map1, map2) // Last write wins

// Shallow copy
copy := record.Clone(m)
```
