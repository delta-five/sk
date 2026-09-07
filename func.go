package sk

// MustMakePtr создаёт указатель *T или вызывает панику при ошибке err.
func MustMakePtr[T any](err error) *T {
	if err != nil {
		panic(err)
	}
	return new(T)
}

// DeRefSafe -
func DeRefSafe[T any](ptr *T) (result T) {
	if ptr != nil {
		return *ptr
	}

	return
}
