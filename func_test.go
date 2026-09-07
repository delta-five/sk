package sk_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/delta-five/sk"
)

func TestDerefOrEmpty(t *testing.T) {
	t.Parallel()

	t.Run("non-nil pointer", func(t *testing.T) {
		t.Parallel()

		val := 42
		assert.Equal(t, 42, sk.DerefOrEmpty(&val))
	})

	t.Run("nil pointer", func(t *testing.T) {
		t.Parallel()

		var p *int
		assert.Zero(t, sk.DerefOrEmpty(p))
	})
}

func ExampleDerefOrEmpty() {
	val := 42
	fmt.Println(sk.DerefOrEmpty(&val))

	var nilPtr *int
	fmt.Println(sk.DerefOrEmpty(nilPtr))
	// Output:
	// 42
	// 0
}
