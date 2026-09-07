package sk

// Option хранит необязательное значение Val с признаком наличия Ok.
type Option[T any] struct {
	Val T
	Ok  bool
}

// Value возвращает хранимое значение. Если значение отсутствует, возвращается нулевое значение типа T.
func (o *Option[T]) Value() T {
	return o.Val
}

// IsValid возвращает true, если значение присутствует.
func (o *Option[T]) IsValid() bool {
	return o.Ok
}

// MakeOption создаёт Option из значения input и признака наличия ok.
func MakeOption[T any](input T, ok bool) Option[T] {
	return Option[T]{Val: input, Ok: ok}
}

// MakeDerefOption создаёт Option из указателя input: при nil — пустой Option, иначе — Option со значением *input.
func MakeDerefOption[T any](input *T) Option[T] {
	if input == nil {
		return Option[T]{}
	}
	return Option[T]{Val: *input, Ok: true}
}

// MakeEmptyOption создаёт пустой Option.
func MakeEmptyOption[T any]() Option[T] {
	return Option[T]{}
}

// OrEmpty возвращает хранимое значение, либо нулевое значение типа T, если значение отсутствует.
func (o *Option[T]) OrEmpty() (result T) {
	if o.Ok {
		return o.Val
	}

	return result
}

// OrValue возвращает хранимое значение, либо переданное значение val, если значение отсутствует.
func (o *Option[T]) OrValue(val T) T {
	if o.Ok {
		return o.Val
	}

	return val
}

// OrElse возвращает хранимое значение, либо результат вызова generator, если значение отсутствует.
// generator вызывается только при отсутствии значения.
func (o *Option[T]) OrElse(generator func() T) T {
	if o.Ok {
		return o.Val
	}

	return generator()
}
