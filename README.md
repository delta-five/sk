# SK - Swiss Knife for Go

The Swiss Army Knife for Go development  — библиотека вспомогательных функций общего назначения.

```bash
go get github.com/delta-five/sk
```

```go
import "github.com/delta-five/sk"
```

## Монада `Result`

Контейнер значения с возможной ошибкой для цепочек преобразований без явной проверки ошибки
на каждом шаге.

### Конструкторы

```go
r := sk.MakeResult(42, nil)        // Result[int] со значением 42
r := sk.MakeResult(0, err)         // Result[int] с ошибкой err
p := sk.MakePtrResult[int](nil)    // Result[*int] с указателем на int
p := sk.MakePtrResult[int](err)    // Result[*int] с ошибкой err
```

### Преобразования

```go
res := sk.MakeResult(21, nil).
	Map(func(n int) string { return strconv.Itoa(n * 2) }) // "42"

res := sk.MakeResult(21, nil).
	MapWithError(func(n int) (string, error) { return parse(n) })

res := sk.MakeResult(21, nil).
	MapWithContext(ctx, func(c context.Context, n int) string { ... })

res := sk.MakeResult(21, nil).
	MapWithContextError(ctx, func(c context.Context, n int) (string, error) { ... })
```

Если исходный `Result` содержит ошибку, `mapper` не вызывается, а ошибка пробрасывается дальше.

Композиция с фабриками срезов:

```go
r := sk.MakeResult([]int{1, 2, 3}, nil).
	Map(sk.NewSliceMapper(func(n int) string { return strconv.Itoa(n) }))
```

### Побочные эффекты

```go
r.Do(func(n int) { fmt.Println(n) })                       // выполняется только при отсутствии ошибки

err := r.DoWithError(func(n int) error { return save(n) }) // возвращает ошибку doer или исходную

r.DoWithContext(ctx, func(c context.Context, n int) { ... })

err := r.DoWithContextError(ctx, func(c context.Context, n int) error { ... })
```

### Извлечение

```go
val, err := r.Unwrap()
```

## Работа со срезами

### `MapSlice`

Применяет `mapper` к каждому элементу среза, возвращает новый срез (порядок и длина сохраняются).

```go
out := sk.MapSlice([]int{1, 2, 3}, func(n int) string { return strconv.Itoa(n * 2) })
// out == []string{"2", "4", "6"}
```

### `MapSliceWithError`

Как `MapSlice`, но прерывается на первой ошибке.

```go
out, err := sk.MapSliceWithError([]int{1, 2, 3}, func(n int) (int, error) {
	if n == 2 {
		return 0, errors.New("boom")
	}
	return n * 2, nil
})
// err != nil, out == nil
```

### `MapSliceWithErrors`

Как `MapSlice`, но не прерывается: обрабатывает все элементы, собирает ошибки
и объединяет их через `errors.Join`.

```go
out, err := sk.MapSliceWithErrors([]int{1, 2, 3, 4}, func(n int) (int, error) {
	if n%2 == 0 {
		return 0, fmt.Errorf("ошибка для %d", n)
	}
	return n * 2, nil
})
// out == nil, err объединяет ошибки для 2 и 4 через errors.Join
```

### Фабрики преобразователей

`NewSliceMapper`, `NewSliceMapperWithError`, `NewSliceMapperWithErrors` возвращают каррированный
преобразователь вида `func([]IN) ([]OUT, error)`, удобный для композиции с `Result`:

```go
mapper := sk.NewSliceMapperWithError(func(n int) (string, error) {
	return strconv.Itoa(n), nil
})
out, err := mapper([]int{1, 2, 3}) // []string{"1", "2", "3"}, nil
```

## Ленивая инициализация

### `NewLazyFunc`

Ленивое вычисление значения без параметров. Генератор вызывается ровно один раз (при первом
вызове), результат кэшируется. Безопасно для конкурентного использования.

```go
lazy := sk.NewLazyFunc(func() *sql.DB {
	db, _ := sql.Open("pgx", dsn)
	return db
})

db := lazy() // первый вызов — инициализация
db = lazy()  // последующие вызовы возвращают кэшированное значение
```

### `NewLazyParamFunc`

Ленивое вычисление значения, параметризованное при первом вызове. Параметр учитывается только
при первом вызове; результат кэшируется и возвращается при всех последующих вызовах.

```go
lazy := sk.NewLazyParamFunc(func(dsn string) *sql.DB {
	db, _ := sql.Open("pgx", dsn)
	return db
})

db := lazy("postgres://...") // первый вызов — инициализация с этим параметром
db = lazy("другой dsn")       // возвращает кэшированное значение, параметр игнорируется
```

### `NewLazyBoundFunc`

Ленивое вычисление значения с заранее фиксированным аргументом. По сути синтаксический сахар
над `NewLazyParamFunc`: аргумент передаётся генератору при первом вызове, результат кэшируется.
Безопасно для конкурентного использования.

```go
lazy := sk.NewLazyBoundFunc("postgres://...", func(dsn string) *sql.DB {
	db, _ := sql.Open("pgx", dsn)
	return db
})

db := lazy() // первый вызов — инициализация с зафиксированным dsn
db = lazy()  // последующие вызовы возвращают кэшированное значение
```
