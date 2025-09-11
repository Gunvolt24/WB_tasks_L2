### Ответ на L2.3

---
Отличная задача с подвохом! 👍

---
Источники информации:
1. https://github.com/golang/go/blob/master/src/internal/abi/iface.go#L14
2. https://github.com/golang/go/blob/253dd08f5df3a45eafc97eec388636fcabfe0174/src/runtime/runtime2.go#L978
3. https://habr.com/ru/articles/449714/
4. https://medium.com/@AlexanderObregon/go-interfaces-behind-the-scenes-f8812706f2c5
5. https://yourbasic.org/golang/gotcha-why-nil-error-not-equal-nil/ 
6. https://medium.com/@ben.meehan_27368/understanding-nil-in-go-interfaces-typed-nil-and-common-pitfalls-6b1154718e00 


---

Программа выведет: 
`<nil>`
`false`

---

#### Что тут происходит:

В функции Foo():

``` go
func Foo() error {
	var err *os.PathError = nil
	return err
}
```

1. `func Foo() error {` - тут важно то, что мы возвращаем `error`. `error` - представляет собой интерфейс:

```go
type error interface {
    Error() string
}
```

1.1. Если копать глубже в представление интерфейса и непустого интерфейса - мы увидим:

```go
// go/src/runtime/runtime2.go

type itab = abi.ITab

type iface struct {
	tab  *itab
	data unsafe.Pointer
}
```

```go
// go/src/internal/abi/iface.go

// The first word of every non-empty interface type contains an *ITab.
// It records the underlying concrete type (Type), the interface type it
// is implementing (Inter), and some ancillary information.
//
// allocated in non-garbage-collected memory
type ITab struct {
	Inter *InterfaceType
	Type  *Type
	...
}

// NonEmptyInterface describes the layout of an interface that contains any methods.
type NonEmptyInterface struct {
	ITab *ITab
	Data unsafe.Pointer
}
```
Здесь можно увидеть, что структура непустого интерфейса хранит в себе **Интерфейсный тип** - `*InterfaceType` (в рамках задачи - интерфейс `error`) и **Конретный тип** - `*Type` (`*os.PathError`). 

2. `var err *os.PathError = nil` - тут объявляем переменную `err` с типом `*os.PathError` и присваиваем ей `nil` значение (!) - тип остается также `*os.PathError` (!). 
**То есть внутри: [тип = `*os.PathError`, значение = `nil`]**
---

В main функции:

```go
func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}
```

1. `err := Foo()` - присваиваем функцию переменной;
2. `fmt.Println(err)` - печатаем **значение** `err` => Вывод:    `[nil]` (печатается nil-значение внутри, но это не значит, что интерфейс nil)
3. `fmt.Println(err == nil)` - вот тут подвох. Как было сказано ранее функция `Foo()` возвращает значение и тип [nil, *os.PathError] и тут мы сравниваем с nil, т.е. с значением и типом [nil, nil]. => Вывод: `false` (интерфейс не пуст: у него есть тип (типизированный nil)).

---
### Сравнение 

По сути, различия только в первых словах структур.

1. Представление пустого интерфейса `(interface{} / any)`:
```go
type EmptyInterface struct {
    Type *Type            // конкретный тип значения
    Data unsafe.Pointer   // адрес данных
}
```
- Здесь первое слово указывает сразу на конкретный тип. Никаких методов.

---
2. Представление непустого интерфейса:

```go
type NonEmptyInterface struct {
    ITab *ITab            // таблица: (интерфейсный тип, конкретный тип, методы)
    Data unsafe.Pointer   // адрес данных
}

type ITab struct {
    Inter *InterfaceType  // сам интерфейсный тип (методный набор)
    Type  *Type           // конкретный тип значения
    Hash  uint32
    Fun   [1]uintptr      // таблица методов 
}
```

- Первое слово - `*ITab`, где хранится связка **интерфейсный тип + конкретный тип** и таблица методов.
---

### Вывод:

Внутренняя структура интерфейса - это пара:

- Динамический тип (у пустого: `*Type`, у непустого: `*ITab` с `Inter+Type`);
- Динамическое значение - `Data`.

Интерфейс считается `nil` только если оба слова равны `nil`: и динамический тип (`Type/ITab`), и динамическое значение (`Data`). Именно поэтому "типизированный nil" (например, `(*T)(nil)` в интерфейсе даёт) `iface != nil`.

Решением такой проблемы является использование Рефлексии.

