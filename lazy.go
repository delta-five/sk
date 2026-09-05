package gsk

import (
	"sync"
)

// ValFunc — функция-генератор значения типа T, не принимающая аргументов.
type ValFunc[T any] func() T

// NewLazyFunc возвращает функцию, которая лениво вычисляет значение типа T при первом вызове.
// Генератор generator вызывается ровно один раз (гарантируется sync.Once), а результат кэшируется
// и возвращается при всех последующих вызовах. Безопасно для конкурентного использования.
func NewLazyFunc[T any](generator ValFunc[T]) ValFunc[T] {
	var (
		once  sync.Once
		value T
	)
	return func() T {
		once.Do(func() {
			value = generator()
		})
		return value
	}
}

// ParamValFunc — функция-генератор значения типа T, принимающая параметр типа P.
type ParamValFunc[P any, T any] func(P) T

// NewLazyParamFunc возвращает функцию, которая лениво вычисляет значение типа T при первом вызове.
// Генератор generator вызывается ровно один раз (гарантируется sync.Once) с параметром, переданным
// при первом вызове, а результат кэшируется и возвращается при всех последующих вызовах независимо
// от значения параметра p. Безопасно для конкурентного использования.
func NewLazyParamFunc[P any, T any](generator ParamValFunc[P, T]) ParamValFunc[P, T] {
	var (
		once  sync.Once
		value T
	)
	return func(p P) T {
		once.Do(func() {
			value = generator(p)
		})
		return value
	}
}

func NewLazyArgFunc[P any, T any](arg P, generator ParamValFunc[P, T]) ValFunc[T] {
	var (
		once  sync.Once
		value T
	)
	return func() T {
		once.Do(func() {
			value = generator(arg)
		})
		return value
	}
}
