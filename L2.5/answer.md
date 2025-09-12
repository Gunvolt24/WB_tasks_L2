### Ответ на L2.5

---
Задача аналогична задаче L2.3 - типизированный nil в непустом интерфейсе.

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
`error`

---

#### На что тут стоит обратить внимание:

В функции main():

``` go
func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

1. `var err error` - это переменная **непустого интерфейса** `error`

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
Здесь можно увидеть, что структура непустого интерфейса хранит в себе **Интерфейсный тип** - `*InterfaceType` и **Конкретный тип** - `*Type`. 

2. `err = test()` - присваиваем переменной `err` функцию `test()`, которая из себя представляет:
```go
func test() *customError {
	// ... do something
	return nil
}
```
- Другими словами, функция `test()` - возвращает **[тип = `*customError`, значение = `nil`]**.

3. Далее сравниваем новую `err` с `nil`:
```go
	if err != nil {
		println("error")
		return
	}
	println("ok")
```

- поскольку переменная `err` представляет собой **[тип = `*customError`, значение = `nil`]** и с ней идет сравнение с nil **[тип = `nil`, значение = `nil`]** => получаем в выводе: `error`, т.к. интерфейс является `nil` **только если оба слова интерфейса ( динамический тип (Type/ITab) и данные (Data)) находятся в состоянии nil**.

---
### Сравнение пустого и непустого интерфейса

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

Интерфейс считается `nil` только если оба слова равны `nil`: и динамический тип (`Type/ITab`), и динамическое значение (`Data`). Именно поэтому "типизированный nil" (например, `(*T)(nil)` в интерфейсе) даёт `iface != nil`.

---

### Вариантами решения такой проблемы являются: 

1. Использование Рефлексии, чтобы увидедть, что внутри `error` хранится `(*customError)(nil)` и избежать непредвиденного результата;
2. Если нужно вернуть `nil`-интерфейс - изменить сигнатуру функции test():
```go
func test() error {
	// do something
	return nil
}
```


