package sk

type Option[T any] struct {
	Val T
	Ok  bool
}

func (o *Option[T]) Value() T {
	return o.Val
}

func (o *Option[T]) IsValid() bool {
	return o.Ok
}

func MakeOption[T any](input T, ok bool) Option[T] {
	return Option[T]{Val: input, Ok: ok}
}

func MakeOptionDeref[T any](input *T) Option[T] {
	if input == nil {
		return Option[T]{}
	}
	return Option[T]{Val: *input, Ok: true}
}

func ToOptionEmpty[T any]() Option[T] {
	return Option[T]{}
}

func (o *Option[T]) OrEmpty() (result T) {
	if o.Ok {
		return o.Val
	}

	return result
}

func (o *Option[T]) OrValue(val T) T {
	if o.Ok {
		return o.Val
	}

	return val
}

func (o *Option[T]) OrMake(generator func() T) T {
	if o.Ok {
		return o.Val
	}

	return generator()
}
