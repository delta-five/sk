package sk

import "errors"

// MapSlice применяет функцию mapper к каждому элементу среза in и возвращает новый срез с результатами.
// Порядок элементов и длина результирующего среза сохраняются.
func MapSlice[IN, OUT any](in []IN, mapper func(IN) OUT) []OUT {
	out := make([]OUT, len(in))
	for i := range in {
		out[i] = mapper(in[i])
	}

	return out
}

// MapSliceWithError применяет функцию mapper к каждому элементу среза in и возвращает новый срез с результатами.
// Если mapper возвращает ошибку для какого-либо элемента, выполнение прерывается и возвращается (nil, err).
func MapSliceWithError[IN, OUT any](in []IN, mapper func(IN) (OUT, error)) ([]OUT, error) {
	out := make([]OUT, len(in))
	for i := range in {
		var err error
		if out[i], err = mapper(in[i]); err != nil {
			return nil, err
		}
	}

	return out, nil
}

// MapSliceWithErrors применяет функцию mapper к каждому элементу среза in и возвращает новый срез с результатами.
// В отличие от MapSliceWithError, не прерывается на первой ошибке: обрабатываются все элементы, а возникшие ошибки
// собираются и возвращаются объединёнными через errors.Join. Если хотя бы один вызов mapper завершился ошибкой,
// возвращается (nil, err), где err — объединённая ошибка.
func MapSliceWithErrors[IN, OUT any](in []IN, mapper func(IN) (OUT, error)) ([]OUT, error) {
	out := make([]OUT, len(in))
	errs := make([]error, 0, len(in))
	for i := range in {
		var err error
		if out[i], err = mapper(in[i]); err != nil {
			errs = append(errs, err)
		}
	}
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}
	return out, nil
}

func NewSliceMapper[IN, OUT any](mapper func(IN) OUT) func([]IN) []OUT {
	return func(in []IN) []OUT {
		return MapSlice[IN, OUT](in, mapper)
	}
}

func NewSliceMapperWithError[IN, OUT any](mapper func(IN) (OUT, error)) func([]IN) ([]OUT, error) {
	return func(in []IN) ([]OUT, error) {
		return MapSliceWithError[IN, OUT](in, mapper)
	}
}

func NewSliceMapperWithErrors[IN, OUT any](mapper func(IN) (OUT, error)) func([]IN) ([]OUT, error) {
	return func(in []IN) ([]OUT, error) {
		return MapSliceWithErrors[IN, OUT](in, mapper)
	}
}
