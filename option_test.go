package sk_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/delta-five/sk"
)

func TestMakeOption(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     int
		ok        bool
		wantVal   int
		wantValid bool
	}{
		{
			name:      "present value",
			input:     42,
			ok:        true,
			wantVal:   42,
			wantValid: true,
		},
		{
			name:    "absent value",
			input:   42,
			ok:      false,
			wantVal: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			o := sk.MakeOption(tt.input, tt.ok)
			assert.Equal(t, tt.wantVal, o.Val)
			assert.Equal(t, tt.wantValid, o.Ok)
		})
	}
}

func TestMakeDerefOption(t *testing.T) {
	t.Parallel()

	t.Run("non-nil pointer", func(t *testing.T) {
		t.Parallel()

		val := 42
		o := sk.MakeDerefOption(&val)
		assert.True(t, o.Ok)
		assert.Equal(t, 42, o.Val)
	})

	t.Run("nil pointer", func(t *testing.T) {
		t.Parallel()

		var p *int
		o := sk.MakeDerefOption(p)
		assert.False(t, o.Ok)
		assert.Zero(t, o.Val)
	})
}

func TestMakeEmptyOption(t *testing.T) {
	t.Parallel()

	o := sk.MakeEmptyOption[int]()
	assert.False(t, o.Ok)
	assert.Zero(t, o.Val)
}

func TestOptionValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		o    sk.Option[int]
		want int
	}{
		{
			name: "present value",
			o:    sk.MakeOption(42, true),
			want: 42,
		},
		{
			name: "absent value",
			o:    sk.MakeEmptyOption[int](),
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.o.Value())
		})
	}
}

func TestOptionIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		o    sk.Option[int]
		want bool
	}{
		{
			name: "present value",
			o:    sk.MakeOption(42, true),
			want: true,
		},
		{
			name: "absent value",
			o:    sk.MakeEmptyOption[int](),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.o.IsValid())
		})
	}
}

func TestOptionOrEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		o    sk.Option[int]
		want int
	}{
		{
			name: "present value",
			o:    sk.MakeOption(42, true),
			want: 42,
		},
		{
			name: "absent value",
			o:    sk.MakeEmptyOption[int](),
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.o.OrEmpty())
		})
	}
}

func TestOptionOrValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		o    sk.Option[int]
		val  int
		want int
	}{
		{
			name: "present value",
			o:    sk.MakeOption(42, true),
			val:  0,
			want: 42,
		},
		{
			name: "absent value",
			o:    sk.MakeEmptyOption[int](),
			val:  7,
			want: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.o.OrValue(tt.val))
		})
	}
}

func TestOptionOrElse(t *testing.T) {
	t.Parallel()

	t.Run("present value does not call generator", func(t *testing.T) {
		t.Parallel()

		called := false
		o := sk.MakeOption(42, true)
		got := o.OrElse(func() int {
			called = true
			return 7
		})

		assert.Equal(t, 42, got)
		assert.False(t, called)
	})

	t.Run("absent value calls generator", func(t *testing.T) {
		t.Parallel()

		called := false
		o := sk.MakeEmptyOption[int]()
		got := o.OrElse(func() int {
			called = true
			return 7
		})

		require.Equal(t, 7, got)
		assert.True(t, called)
	})
}

func ExampleMakeOption() {
	opt := sk.MakeOption(42, true)
	fmt.Println(opt.Value(), opt.IsValid())

	empty := sk.MakeOption(0, false)
	fmt.Println(empty.Value(), empty.IsValid())
	// Output:
	// 42 true
	// 0 false
}

func ExampleMakeDerefOption() {
	val := 42
	opt := sk.MakeDerefOption(&val)
	fmt.Println(opt.Value(), opt.IsValid())

	var nilPtr *int
	empty := sk.MakeDerefOption(nilPtr)
	fmt.Println(empty.Value(), empty.IsValid())
	// Output:
	// 42 true
	// 0 false
}

func ExampleOption_OrValue() {
	opt := sk.MakeOption(42, true)
	fmt.Println(opt.OrValue(7))

	empty := sk.MakeEmptyOption[int]()
	fmt.Println(empty.OrValue(7))
	// Output:
	// 42
	// 7
}
