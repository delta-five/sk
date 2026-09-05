package sk_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/delta-five/sk"
)

func TestMapSlice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		in     []int
		mapper func(int) string
		want   []string
	}{
		{
			name:   "empty slice",
			in:     []int{},
			mapper: strconv.Itoa,
			want:   []string{},
		},
		{
			name:   "non-empty slice",
			in:     []int{1, 2, 3},
			mapper: func(n int) string { return strconv.Itoa(n * 2) },
			want:   []string{"2", "4", "6"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := sk.MapSlice(tt.in, tt.mapper)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMapSliceError(t *testing.T) {
	t.Parallel()

	errBoom := errors.New("boom")

	tests := []struct {
		name    string
		in      []int
		mapper  func(int) (int, error)
		want    []int
		wantErr error
	}{
		{
			name:   "empty slice",
			in:     []int{},
			mapper: func(n int) (int, error) { return n, nil },
			want:   []int{},
		},
		{
			name:   "success",
			in:     []int{1, 2, 3},
			mapper: func(n int) (int, error) { return n * 2, nil },
			want:   []int{2, 4, 6},
		},
		{
			name: "error",
			in:   []int{1, 2, 3},
			mapper: func(n int) (int, error) {
				if n == 2 {
					return 0, errBoom
				}
				return n, nil
			},
			wantErr: errBoom,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := sk.MapSliceWithError(tt.in, tt.mapper)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMapSliceErrors(t *testing.T) {
	t.Parallel()

	errBoom := errors.New("boom")
	errBang := errors.New("bang")

	tests := []struct {
		name       string
		in         []int
		mapper     func(int) (int, error)
		want       []int
		wantErrors []error
	}{
		{
			name:   "empty slice",
			in:     []int{},
			mapper: func(n int) (int, error) { return n, nil },
			want:   []int{},
		},
		{
			name:   "success",
			in:     []int{1, 2, 3},
			mapper: func(n int) (int, error) { return n * 2, nil },
			want:   []int{2, 4, 6},
		},
		{
			name: "single error",
			in:   []int{1, 2, 3},
			mapper: func(n int) (int, error) {
				if n == 2 {
					return 0, errBoom
				}
				return n, nil
			},
			wantErrors: []error{errBoom},
		},
		{
			name: "multiple errors",
			in:   []int{1, 2, 3},
			mapper: func(n int) (int, error) {
				switch n {
				case 1:
					return 0, errBoom
				case 2:
					return 0, errBang
				default:
					return n, nil
				}
			},
			wantErrors: []error{errBoom, errBang},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := sk.MapSliceWithErrors(tt.in, tt.mapper)
			if len(tt.wantErrors) > 0 {
				require.Error(t, err)
				assert.Nil(t, got)
				for _, wantErr := range tt.wantErrors {
					require.ErrorIs(t, err, wantErr)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
