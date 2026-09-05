package gsk_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/delta-five/gsk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLazyFunc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want int
	}{
		{
			name: "computes value once",
			want: 42,
		},
		{
			name: "returns cached value",
			want: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var calls atomic.Int32
			lazy := gsk.NewLazyFunc(func() int {
				calls.Add(1)
				return tt.want
			})

			first := lazy()
			second := lazy()
			third := lazy()

			assert.Equal(t, tt.want, first)
			assert.Equal(t, tt.want, second)
			assert.Equal(t, tt.want, third)
			assert.Equal(t, int32(1), calls.Load())
		})
	}
}

func TestNewLazyFuncConcurrent(t *testing.T) {
	t.Parallel()

	var calls atomic.Int32
	lazy := gsk.NewLazyFunc(func() int {
		calls.Add(1)
		return 42
	})

	const goroutines = 100
	var wg sync.WaitGroup
	results := make(chan int, goroutines)

	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			results <- lazy()
		}()
	}
	wg.Wait()
	close(results)

	for got := range results {
		assert.Equal(t, 42, got)
	}
	assert.Equal(t, int32(1), calls.Load())
}

func TestNewLazyParamFunc(t *testing.T) {
	t.Parallel()

	t.Run("computes value once with first param", func(t *testing.T) {
		t.Parallel()

		var calls atomic.Int32
		var captured atomic.Value
		lazy := gsk.NewLazyParamFunc(func(p int) int {
			calls.Add(1)
			captured.Store(p)
			return p * 2
		})

		first := lazy(5)
		second := lazy(10)
		third := lazy(15)

		assert.Equal(t, 10, first)
		assert.Equal(t, 10, second)
		assert.Equal(t, 10, third)
		assert.Equal(t, int32(1), calls.Load())
		assert.Equal(t, 5, captured.Load())
	})
}

func TestNewLazyParamFuncConcurrent(t *testing.T) {
	t.Parallel()

	var calls atomic.Int32
	var captured atomic.Value
	lazy := gsk.NewLazyParamFunc(func(p int) int {
		calls.Add(1)
		captured.Store(p)
		return p * 2
	})

	const goroutines = 100
	var wg sync.WaitGroup
	results := make(chan int, goroutines)

	wg.Add(goroutines)
	for i := range goroutines {
		go func(p int) {
			defer wg.Done()
			results <- lazy(p)
		}(i)
	}
	wg.Wait()
	close(results)

	for got := range results {
		require.Equal(t, captured.Load(), got/2)
	}
	assert.Equal(t, int32(1), calls.Load())
}

func ExampleNewLazyFunc() {
	counter := 0
	lazy := gsk.NewLazyFunc(func() int {
		counter++
		return counter * 10
	})

	fmt.Println(lazy())
	fmt.Println(lazy())
	// Output:
	// 10
	// 10
}

func ExampleNewLazyParamFunc() {
	calls := 0
	lazy := gsk.NewLazyParamFunc(func(p int) int {
		calls++
		return p + calls
	})

	fmt.Println(lazy(1))
	fmt.Println(lazy(100))
	// Output:
	// 2
	// 2
}
