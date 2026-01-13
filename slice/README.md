# Slice Package

The `slice` package provides utility functions for working with Go slices, featuring a powerful `Collect` function for complex iteration, filtering, and error handling scenarios.

## Features

- **Ergonomic Flow Control**: Use `Stop()`, `Continue()`, and `GetValue()` with immediate effect.
- **Error Accumulation**: Automatically collects errors encountered during iteration.
- **Type Safety**: Generic implementation supports any types.

## Usage

### Basic Collection

```go
input := []int{1, 2, 3, 4, 5}
result, err := slice.Collect(input, func(c slice.CollectorContext[int, int]) {
    _, val := c.CurrentElem()
    if val % 2 == 0 {
        c.SetValue(val * 2)
    }
})
// result: [4, 8]
```

### Flow Control with Errors

```go
result, err := slice.Collect(input, func(c slice.CollectorContext[int, int]) {
    _, val := c.CurrentElem()
    
    if isFatal(val) {
        c.Stop(errors.New("fatal error")) // Stops immediately
    }
    
    if shouldSkip(val) {
        c.Continue() // Skips immediately
    }
    
    c.SetValue(val)
})
```

## Performance & Use Cases

### Benchmark Results

| Implementation | Time/Op | Allocations | Rank |
| :--- | :--- | :--- | :--- |
| **Native Loop** | `~1,767 ns/op` | 10 allocs/op | 1 |
| **Collect Fn** | `~9,320 ns/op` | 16 allocs/op | 2 |

*Benchmarks run on Apple M2 Pro processing 1,000 items.*

### Recommendations

**✅ Good For:**
- Business logic pipelines
- Complex filtering requirements
- Code that benefits from consolidated error handling
- CLI tools and data processing scripts

**❌ Avoid For:**
- High-performance "hot paths" (e.g., game loops, trading engines)
- Latency-critical systems where microseconds matter
- Operations performing simple iterating without complex flow logic

The ~5x overhead comes from the abstraction layer (closures, interface calls, defer/recover) which provides the improved ergonomics.
