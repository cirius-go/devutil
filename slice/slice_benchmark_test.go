package slice

import "testing"

func BenchmarkCollect(b *testing.B) {
	input := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		input[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Collect(input, func(c CollectorContext[int, int]) {
			_, val := c.CurrentElem()
			if val%2 == 0 {
				c.SetValue(val)
			}
		})
	}
}

func BenchmarkCollect_SimpleCopy(b *testing.B) {
	input := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		input[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Collect(input, func(c CollectorContext[int, int]) {
			_, val := c.CurrentElem()
			c.SetValue(val)
		})
	}
}

func BenchmarkNativeLoop(b *testing.B) {
	input := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		input[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result []int
		for _, val := range input {
			if val%2 == 0 {
				result = append(result, val)
			}
		}
	}
}
