# DevUtil

A collection of utility packages for Go, providing ergonomic and type-safe helpers for common operations on slices and maps.

## Packages

### [slice](./slice)

Utilities for working with slices, including functional programming patterns and advanced collection operations.

**Key Features:**
- `Collect`: Advanced iteration with immediate control flow (`Stop`/`Continue`) and rich error handling
- `Filter`, `Map`, `Reduce`: Standard functional operations
- `Every`, `Some`, `Find`, `Contains`: Predicate-based queries
- `SliceError`: Detailed error tracking with element context (index, value)

**Example:**
```go
import "github.com/cirius-go/devutil/slice"

result, err := slice.Collect(input, func(c slice.CollectorContext[int, int]) {
    _, val := c.CurrentElem()
    if shouldSkip(val) {
        c.Continue() // Skip immediately
    }
    if isFatal(val) {
        c.Stop(err) // Stop immediately with error
    }
    c.SetValue(val * 2)
})
```

[Read more →](./slice/README.md)

### [record](./record)

Utilities for working with maps (records/associative arrays), providing common transformation and access patterns.

**Key Features:**
- `Keys`, `Values`: Extract to slices
- `SortedKeys`, `SortedValues`: Ordered extraction
- `Clone`, `Merge`: Map operations
- `Filter`, `MapValues`: Transformations
- `ToSet`: Convert slice to set
- `Associate`: Build map from slice with transform

**Example:**
```go
import "github.com/cirius-go/devutil/record"

// Create a set from slice
set := record.ToSet([]string{"a", "b", "a"})
// map[string]struct{}{"a":{}, "b":{}}

// Transform slice to map
lengths := record.Associate(words, func(w string) (string, int) {
    return w, len(w)
})
```

[Read more →](./record/README.md)

## Installation

```bash
go get github.com/cirius-go/devutil
```

## Requirements

- Go 1.21+ (uses generics and standard library features like `cmp.Ordered`)

## Performance

- **slice.Collect**: ~5x slower than native loops due to abstraction overhead. Best for business logic where ergonomics matter.
- **Utility functions**: Rank 1 performance (native loop speed). Use for simple operations.

## License

[Add your license here]
