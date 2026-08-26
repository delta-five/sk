package ga_test

import (
	"context"
	"errors"
	"maps"
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ga "github.com/MirrorRu/ga"
)

func TestMakeResult(t *testing.T) {
	t.Parallel()

	t.Run("without error", func(t *testing.T) {
		t.Parallel()

		r := ga.MakeResult(42, nil)
		val, err := r.Unwrap()
		assert.Equal(t, 42, val)
		assert.NoError(t, err)
	})

	t.Run("with error", func(t *testing.T) {
		t.Parallel()

		r := ga.MakeResult(0, assert.AnError)
		val, err := r.Unwrap()
		assert.Zero(t, val)
		assert.ErrorIs(t, err, assert.AnError)
	})
}

func TestUnwrap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		r       ga.Result[int]
		want    int
		wantErr error
	}{
		{
			name: "with value",
			r:    ga.MakeResult(42, nil),
			want: 42,
		},
		{
			name:    "with error",
			r:       ga.MakeResult(0, assert.AnError),
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			val, err := tt.r.Unwrap()
			assert.Equal(t, tt.want, val)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestResultMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		r       ga.Result[int]
		mapper  func(int) string
		want    string
		wantErr error
	}{
		{
			name:   "success",
			r:      ga.MakeResult(21, nil),
			mapper: func(n int) string { return strconv.Itoa(n * 2) },
			want:   "42",
		},
		{
			name:    "error",
			r:       ga.MakeResult(21, assert.AnError),
			mapper:  func(n int) string { return "should not run" },
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.r.Map(tt.mapper)
			val, err := got.Unwrap()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Zero(t, val)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, val)
		})
	}
}

func TestResultMapSlice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		r       ga.Result[[]int]
		mapper  func(int) string
		want    []string
		wantErr error
	}{
		{
			name:   "success",
			r:      ga.MakeResult([]int{1, 2, 3}, nil),
			mapper: func(n int) string { return strconv.Itoa(n * n) },
			want:   []string{"1", "4", "9"},
		},
		{
			name:    "error",
			r:       ga.MakeResult([]int(nil), assert.AnError),
			mapper:  func(n int) string { return "should not run" },
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.r.Map(ga.NewSliceMapper(tt.mapper))
			val, err := got.Unwrap()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Zero(t, val)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, val)
		})
	}
}

func TestResultMapWithError(t *testing.T) {
	t.Parallel()

	mapperErr := errors.New("mapper error")

	tests := []struct {
		name    string
		r       ga.Result[int]
		mapper  func(int) (string, error)
		want    string
		wantErr error
	}{
		{
			name:   "success",
			r:      ga.MakeResult(21, nil),
			mapper: func(n int) (string, error) { return strconv.Itoa(n * 2), nil },
			want:   "42",
		},
		{
			name:    "mapper error",
			r:       ga.MakeResult(21, nil),
			mapper:  func(n int) (string, error) { return "", mapperErr },
			wantErr: mapperErr,
		},
		{
			name:    "result error",
			r:       ga.MakeResult(21, assert.AnError),
			mapper:  func(n int) (string, error) { return strconv.Itoa(n), nil },
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.r.MapWithError(tt.mapper)
			val, err := got.Unwrap()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Zero(t, val)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, val)
		})
	}
}

func TestResultMapSliceWithError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		r       ga.Result[[]int]
		mapper  func(int) (string, error)
		want    []string
		wantErr error
	}{
		{
			name:   "success",
			r:      ga.MakeResult([]int{1, 2, 3}, nil),
			mapper: func(n int) (string, error) { return strconv.Itoa(n), nil },
			want:   []string{"1", "2", "3"},
		},
		{
			name:    "mapper error",
			r:       ga.MakeResult([]int{1, 2, 3}, nil),
			mapper:  func(n int) (string, error) { return "", assert.AnError },
			wantErr: assert.AnError,
		},
		{
			name:    "result error",
			r:       ga.MakeResult([]int(nil), assert.AnError),
			mapper:  func(n int) (string, error) { return strconv.Itoa(n), nil },
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.r.MapWithError(ga.NewSliceMapperWithError(tt.mapper))
			val, err := got.Unwrap()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Zero(t, val)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, val)
		})
	}
}

func TestResultMapSliceWithErrors(t *testing.T) {
	t.Parallel()

	errMap := map[int]error{
		1: errors.New("mapper error 1"),
		2: errors.New("mapper error 2"),
		3: nil,
	}
	valMap := map[int]string{
		1: "1",
		2: "2",
		3: "3",
	}

	require.Equal(t, slices.Sorted[int](maps.Keys(valMap)), slices.Sorted[int](maps.Keys(errMap)))

	tests := []struct {
		name    string
		r       ga.Result[[]int]
		mapper  func(int) (string, error)
		want    []string
		wantErr []error
	}{
		{
			name:   "success",
			r:      ga.MakeResult(slices.Collect(maps.Keys(errMap)), nil),
			mapper: func(n int) (string, error) { return strconv.Itoa(n), nil },
			want:   slices.Sorted(maps.Values(valMap)),
		},
		{
			name: "mapper error",
			r:    ga.MakeResult([]int{1, 2, 3}, nil),
			mapper: func(n int) (string, error) {
				val, valOk := valMap[n]
				require.True(t, valOk)
				err, errOk := errMap[n]
				require.True(t, errOk)
				return val, err
			},
			wantErr: slices.Collect(maps.Values(errMap)),
		},
		{
			name:    "result error",
			r:       ga.MakeResult([]int(nil), assert.AnError),
			mapper:  func(n int) (string, error) { return strconv.Itoa(n), nil },
			wantErr: []error{assert.AnError},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.r.MapWithError(ga.NewSliceMapperWithErrors(tt.mapper))
			val, err := got.Unwrap()
			if tt.wantErr != nil {
				for _, e := range tt.wantErr {
					if e != nil {
						assert.ErrorIs(t, err, e)
					}
				}
				assert.Zero(t, val)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, val)
		})
	}
}

func TestResultMapWithContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		r       ga.Result[int]
		ctx     context.Context
		mapper  func(context.Context, int) string
		want    string
		wantErr error
	}{
		{
			name:   "success",
			r:      ga.MakeResult(21, nil),
			ctx:    context.Background(),
			mapper: func(ctx context.Context, n int) string { return strconv.Itoa(n * 2) },
			want:   "42",
		},
		{
			name:    "error",
			r:       ga.MakeResult(21, assert.AnError),
			ctx:     context.Background(),
			mapper:  func(ctx context.Context, n int) string { return "should not run" },
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.r.MapWithContext(tt.ctx, tt.mapper)
			val, err := got.Unwrap()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Zero(t, val)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, val)
		})
	}
}

func TestResultMapWithContextError(t *testing.T) {
	t.Parallel()

	mapperErr := errors.New("mapper error")

	tests := []struct {
		name    string
		r       ga.Result[int]
		ctx     context.Context
		mapper  func(context.Context, int) (string, error)
		want    string
		wantErr error
	}{
		{
			name:   "success",
			r:      ga.MakeResult(21, nil),
			ctx:    context.Background(),
			mapper: func(ctx context.Context, n int) (string, error) { return strconv.Itoa(n * 2), nil },
			want:   "42",
		},
		{
			name:    "mapper error",
			r:       ga.MakeResult(21, nil),
			ctx:     context.Background(),
			mapper:  func(ctx context.Context, n int) (string, error) { return "", mapperErr },
			wantErr: mapperErr,
		},
		{
			name:    "result error",
			r:       ga.MakeResult(21, assert.AnError),
			ctx:     context.Background(),
			mapper:  func(ctx context.Context, n int) (string, error) { return strconv.Itoa(n), nil },
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.r.MapWithContextError(tt.ctx, tt.mapper)
			val, err := got.Unwrap()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Zero(t, val)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, val)
		})
	}
}

func TestResultDo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		r    ga.Result[int]
		want []int
	}{
		{
			name: "success",
			r:    ga.MakeResult(42, nil),
			want: []int{42},
		},
		{
			name: "error",
			r:    ga.MakeResult(0, assert.AnError),
			want: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			captured := make([]int, 0)
			tt.r.Do(func(n int) { captured = append(captured, n) })
			assert.Equal(t, tt.want, captured)
		})
	}
}

func TestResultDoWithError(t *testing.T) {
	t.Parallel()

	doerErr := errors.New("doer error")

	tests := []struct {
		name    string
		r       ga.Result[int]
		doer    func(int) error
		wantErr error
	}{
		{
			name: "success returns nil",
			r:    ga.MakeResult(42, nil),
			doer: func(n int) error { return nil },
		},
		{
			name:    "success returns error",
			r:       ga.MakeResult(42, nil),
			doer:    func(n int) error { return doerErr },
			wantErr: doerErr,
		},
		{
			name:    "result error",
			r:       ga.MakeResult(0, assert.AnError),
			doer:    func(n int) error { return nil },
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.r.DoWithError(tt.doer)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestResultDoWithContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		r    ga.Result[int]
		ctx  context.Context
		want []int
	}{
		{
			name: "success",
			r:    ga.MakeResult(42, nil),
			ctx:  context.Background(),
			want: []int{42},
		},
		{
			name: "error",
			r:    ga.MakeResult(0, assert.AnError),
			ctx:  context.Background(),
			want: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			captured := make([]int, 0)
			tt.r.DoWithContext(tt.ctx, func(ctx context.Context, n int) {
				captured = append(captured, n)
			})
			assert.Equal(t, tt.want, captured)
		})
	}
}

func TestResultDoWithContextError(t *testing.T) {
	t.Parallel()

	doerErr := errors.New("doer error")

	tests := []struct {
		name    string
		r       ga.Result[int]
		ctx     context.Context
		doer    func(context.Context, int) error
		wantErr error
	}{
		{
			name: "success returns nil",
			r:    ga.MakeResult(42, nil),
			ctx:  context.Background(),
			doer: func(ctx context.Context, n int) error { return nil },
		},
		{
			name:    "success returns error",
			r:       ga.MakeResult(42, nil),
			ctx:     context.Background(),
			doer:    func(ctx context.Context, n int) error { return doerErr },
			wantErr: doerErr,
		},
		{
			name:    "result error",
			r:       ga.MakeResult(0, assert.AnError),
			ctx:     context.Background(),
			doer:    func(ctx context.Context, n int) error { return nil },
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.r.DoWithContextError(tt.ctx, tt.doer)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
