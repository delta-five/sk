package sk

// MustDo проверяет, не является ли аргумент ошибкой, и в противном случае выдает ошибку.
// Обычно аргумент передается напрямую как результат вызова другого метода.
func MustDo(err error) {
	if err != nil {
		panic(err)
	}
}

// MustMake создаёт T из значения val или вызывает панику при ошибке err.
func MustMake[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

// MustAllocate создаёт указатель *T или вызывает панику при ошибке err.
func MustAllocate[T any](err error) *T {
	if err != nil {
		panic(err)
	}
	return new(T)
}

// DerefOrEmpty возвращает значение, на которое указывает ptr, либо нулевое значение типа T, если ptr == nil.
func DerefOrEmpty[T any](ptr *T) (result T) {
	if ptr != nil {
		return *ptr
	}

	return
}
