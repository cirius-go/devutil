package slice

import (
	"sync"
)

// PipeFn defines a function type that processes an item of type In and returns
// an item of the same type.
type PipeFn[In any] func(item In) In

type controlSignal int

const (
	sigExclude controlSignal = iota // default, do nothing (should not happen explicitly mostly)
	sigStop
	sigContinue
)

// collectorContextImpl is a concrete implementation of CollectorContext.
type collectorContextImpl[In, Out any] struct {
	// global state
	copiedOnce   *sync.Once
	copied       []In
	sliceGetter  func() []In
	elemGetter   func(index int) In
	resultGetter func() []Out
	// signals
	continued      bool
	stopped        bool
	errOnStopped   []*ElemError[In]
	errOnContinued []*ElemError[In]
	// current state
	currentIndex  int
	currentResult []Out
	currentValue  Out
	hasValue      bool
}

// Slice implements the Slice method of CollectorContext.
func (c *collectorContextImpl[In, Out]) Slice() []In {
	return c.sliceGetter()
}

// CurrentElem implements the CurrentElem method of CollectorContext.
func (c *collectorContextImpl[In, Out]) CurrentElem() (int, In) {
	return c.currentIndex, c.elemGetter(c.currentIndex)
}

// SetValue implements the SetValue method of CollectorContext.
func (c *collectorContextImpl[In, Out]) SetValue(value Out) {
	c.currentValue = value
	c.hasValue = true
}

// CurrentResult implements the CurrentResult method of CollectorContext.
func (c *collectorContextImpl[In, Out]) CurrentResult() []Out {
	return c.resultGetter()
}

// Continue implements the Continue method of CollectorContext.
func (c *collectorContextImpl[In, Out]) Continue(errs ...error) {
	for _, err := range errs {
		if err != nil {
			c.errOnContinued = append(c.errOnContinued, &ElemError[In]{
				Index: c.currentIndex,
				Value: c.elemGetter(c.currentIndex),
				Err:   err,
			})
		}
	}
	c.continued = true
	panic(sigContinue)
}

// Stop implements the Stop method of CollectorContext.
func (c *collectorContextImpl[In, Out]) Stop(errs ...error) {
	for _, err := range errs {
		if err != nil {
			c.errOnStopped = append(c.errOnStopped, &ElemError[In]{
				Index: c.currentIndex,
				Value: c.elemGetter(c.currentIndex),
				Err:   err,
			})
		}
	}
	c.stopped = true
	panic(sigStop)
}

// CollectorContext is a context interface for collection operations.
type CollectorContext[In, Out any] interface {
	// Slice returns a copy of the original input slice.
	Slice() []In
	// CurrentElem returns the index and value of the current element being processed.
	CurrentElem() (int, In)
	// SetValue sets the value to be added to the result.
	SetValue(value Out)
	// CurrentResult returns a copy of the current result slice.
	CurrentResult() []Out
	// Continue signals to skip adding the current element to the result.
	Continue(errs ...error)
	// Stop signals to terminate the collection process immediately.
	Stop(errs ...error)
}

// Collect applies a collection operation on the input slice based on the provided context,
// and returns an error if the handler fails.
func Collect[In, Out any](input []In, handler func(c CollectorContext[In, Out])) ([]Out, error) {
	var (
		result []Out
		errs   SliceError[In]
	)
	if len(input) == 0 || handler == nil {
		return result, nil
	}

	// construct context.
	c := &collectorContextImpl[In, Out]{
		copiedOnce:  &sync.Once{},
		copied:      nil,
		sliceGetter: nil,
		elemGetter:  nil,

		continued:      false,
		stopped:        false,
		errOnStopped:   nil,
		errOnContinued: nil,

		currentIndex:  0,
		currentResult: nil,
		hasValue:      false,
	}
	c.sliceGetter = func() []In {
		if c.copied == nil {
			c.copiedOnce.Do(func() {
				copied := make([]In, len(input))
				copy(copied, input)
				c.copied = copied
			})
		}
		return c.copied
	}
	c.elemGetter = func(index int) In {
		return input[index]
	}
	c.resultGetter = func() []Out {
		if len(result) == 0 {
			return nil
		}
		res := make([]Out, len(result))
		copy(res, result)
		return res
	}

	for i := range input {
		c.stopped = false
		c.continued = false
		c.errOnStopped = nil
		c.errOnContinued = nil
		c.hasValue = false
		// zero out currentValue? Not strictly necessary if hasValue checks cover it,
		// but good for GC if Out is a pointer/large struct.
		// For generic 'Out', we can't easily set to zero value without reflection or declaring var.
		// We'll rely on hasValue.

		c.currentIndex = i

		func() {
			defer func() {
				if r := recover(); r != nil {
					if s, ok := r.(controlSignal); ok {
						if s == sigStop {
							c.stopped = true
							return
						}
						if s == sigContinue {
							c.continued = true
							return
						}
					}
					panic(r)
				}
			}()
			handler(c)
		}()

		if c.stopped {
			if len(c.errOnStopped) > 0 {
				errs = append(errs, c.errOnStopped...)
			}
			break
		}
		if c.continued {
			if len(c.errOnContinued) > 0 {
				errs = append(errs, c.errOnContinued...)
			}
			continue
		}

		if c.hasValue {
			result = append(result, c.currentValue)
		}
	}
	if len(errs) == 0 {
		return result, nil
	}
	return result, errs
}

// Filter applies a filtering operation on the input slice based on the provided predicate function.
func Filter[In any](input []In, predicate func(item In) bool) []In {
	if len(input) == 0 || predicate == nil {
		return input
	}
	res, _ := Collect(input, func(c CollectorContext[In, In]) {
		_, val := c.CurrentElem()
		if !predicate(val) {
			c.Continue()
		}
		c.SetValue(val)
	})
	return res
}

// Reduce applies a reduction operation on the input slice based on the provided reducer function.
func Reduce[In, Out any](input []In, reducer func(Out, In) Out, initial Out) Out {
	if len(input) == 0 || reducer == nil {
		return initial
	}
	var acc Out = initial
	for _, item := range input {
		acc = reducer(acc, item)
	}
	return acc
}

// Every returns true if all elements in the slice satisfy the predicate.
// Returns true for empty slices (vacuously true).
func Every[In any](input []In, predicate func(item In) bool) bool {
	if predicate == nil {
		return true
	}
	for _, item := range input {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// Some returns true if at least one element in the slice satisfies the predicate.
// Returns false for empty slices.
func Some[In any](input []In, predicate func(item In) bool) bool {
	if len(input) == 0 || predicate == nil {
		return false
	}
	for _, item := range input {
		if predicate(item) {
			return true
		}
	}
	return false
}

// Map applies a transformation function to each element of the slice and returns a new slice.
func Map[In, Out any](input []In, mapper func(item In) Out) []Out {
	if len(input) == 0 || mapper == nil {
		return nil
	}
	result := make([]Out, len(input))
	for i, item := range input {
		result[i] = mapper(item)
	}
	return result
}

// Find returns the first element that satisfies the predicate and true.
// If no element matches, it returns the zero value and false.
func Find[In any](input []In, predicate func(item In) bool) (In, bool) {
	var zero In
	if len(input) == 0 || predicate == nil {
		return zero, false
	}
	for _, item := range input {
		if predicate(item) {
			return item, true
		}
	}
	return zero, false
}

// Contains returns true if the slice contains the target element.
// In must be comparable.
func Contains[In comparable](input []In, target In) bool {
	if len(input) == 0 {
		return false
	}
	for _, item := range input {
		if item == target {
			return true
		}
	}
	return false
}

// Chunk splits a slice into chunks of the specified size.
// If the slice cannot be split evenly, the last chunk will contain the remaining elements.
// Returns nil if the input slice is nil.
// If size is <= 0, it defaults to 1.
func Chunk[In any](input []In, size int) [][]In {
	if input == nil {
		return nil
	}
	if len(input) == 0 {
		return make([][]In, 0)
	}
	if size <= 0 {
		size = 1
	}

	chunks := make([][]In, 0, (len(input)+size-1)/size)
	for size < len(input) {
		input, chunks = input[size:], append(chunks, input[0:size:size])
	}
	chunks = append(chunks, input)
	return chunks
}

// ForEachChunk splits the slice into chunks and processes them using the handler.
// The concurrency parameter controls the number of concurrent handlers.
// If concurrency <= 1, chunks are processed sequentially.
// If any handler returns an error, the function returns the first error encountered.
// Note: When running concurrently, the order of execution is not guaranteed,
// and it will wait for all started goroutines to finish even if one fails.
func ForEachChunk[In any](input []In, chunkSize int, concurrency int, handler func(chunk []In) error) error {
	if len(input) == 0 {
		return nil
	}
	chunks := Chunk(input, chunkSize)

	if concurrency <= 1 {
		for _, chunk := range chunks {
			if err := handler(chunk); err != nil {
				return err
			}
		}
		return nil
	}

	var (
		wg      sync.WaitGroup
		errChan = make(chan error, len(chunks))
		sem     = make(chan struct{}, concurrency)
	)

	for _, chunk := range chunks {
		sem <- struct{}{} // Acquire token
		wg.Add(1)
		go func(c []In) {
			defer wg.Done()
			defer func() { <-sem }() // Release token
			if err := handler(c); err != nil {
				errChan <- err
			}
		}(chunk)
	}

	wg.Wait()
	close(errChan)

	// Return the first error if any
	if len(errChan) > 0 {
		return <-errChan
	}
	return nil
}

// Flatten flattens a slice of slices into a single slice.
// It pre-allocates the result slice to minimize allocations.
// Returns nil if input is nil.
func Flatten[In any](input [][]In) []In {
	if input == nil {
		return nil
	}
	if len(input) == 0 {
		return []In{}
	}

	totalLen := 0
	for _, s := range input {
		totalLen += len(s)
	}

	result := make([]In, 0, totalLen)
	for _, s := range input {
		result = append(result, s...)
	}
	return result
}
