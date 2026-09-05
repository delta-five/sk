package gsk

import (
	"context"
)

// Result хранит значение val и возможную ошибку err, позволяя строить цепочки преобразований
// без явной проверки ошибки на каждом шаге.
type Result[IN any] struct {
	val IN
	err error
}

// MakeResult создаёт Result из значения val и ошибки err.
func MakeResult[T any](val T, err error) Result[T] {
	return Result[T]{val: val, err: err}
}

// MakePtrResult создаёт Result с указателем на переменную с типом параметра при пустом аргументе-ошибки.
func MakePtrResult[T any](err error) Result[*T] {
	if err == nil {
		return Result[*T]{val: new(T)}
	}
	return Result[*T]{err: err}
}

// Unwrap возвращает хранимое значение и ошибку.
func (r *Result[IN]) Unwrap() (IN, error) {
	return r.val, r.err
}

// Map применяет mapper к значению, если ошибки нет. При наличии ошибки возвращает Result с той же ошибкой.
func (r *Result[IN]) Map[OUT any](mapper func(IN) OUT) Result[OUT] {
	if r.err != nil {
		return Result[OUT]{err: r.err}
	}

	return Result[OUT]{val: mapper(r.val)}
}

// MapWithError применяет mapper к значению, если ошибки нет, и возвращает результат вызова mapper.
// При наличии исходной ошибки возвращает Result с той же ошибкой, не вызывая mapper.
func (r *Result[IN]) MapWithError[OUT any](mapper func(IN) (OUT, error)) Result[OUT] {
	if r.err != nil {
		return Result[OUT]{err: r.err}
	}

	return MakeResult(mapper(r.val))
}

// MapWithContext применяет mapper к значению и контексту ctx, если ошибки нет.
// При наличии ошибки возвращает Result с той же ошибкой, не вызывая mapper.
func (r *Result[IN]) MapWithContext[OUT any](ctx context.Context, mapper func(context.Context, IN) OUT) Result[OUT] {
	if r.err != nil {
		return Result[OUT]{err: r.err}
	}

	return Result[OUT]{val: mapper(ctx, r.val)}
}

// MapWithContextError применяет mapper к значению и контексту ctx, если ошибки нет, и возвращает результат вызова mapper.
// При наличии исходной ошибки возвращает Result с той же ошибкой, не вызывая mapper.
func (r *Result[IN]) MapWithContextError[OUT any](ctx context.Context, mapper func(context.Context, IN) (OUT, error)) Result[OUT] {
	if r.err != nil {
		return Result[OUT]{err: r.err}
	}

	return MakeResult(mapper(ctx, r.val))
}

// Do вызывает doer со значением, если ошибки нет. При наличии ошибки ничего не делает.
func (r *Result[IN]) Do(doer func(IN)) {
	if r.err == nil {
		doer(r.val)
	}
}

// DoWithError вызывает doer со значением, если ошибки нет, и возвращает результат его вызова.
// При наличии исходной ошибки возвращает эту ошибку, не вызывая doer.
func (r *Result[IN]) DoWithError(doer func(IN) error) error {
	if r.err == nil {
		return doer(r.val)
	}

	return r.err
}

// DoWithContext вызывает doer со значением и контекстом ctx, если ошибки нет.
// При наличии ошибки ничего не делает.
func (r *Result[IN]) DoWithContext(ctx context.Context, doer func(context.Context, IN)) {
	if r.err == nil {
		doer(ctx, r.val)
	}
}

// DoWithContextError вызывает doer со значением и контекстом ctx, если ошибки нет, и возвращает результат его вызова.
// При наличии исходной ошибки возвращает эту ошибку, не вызывая doer.
func (r *Result[IN]) DoWithContextError(ctx context.Context, doer func(context.Context, IN) error) error {
	if r.err == nil {
		return doer(ctx, r.val)
	}

	return r.err
}
