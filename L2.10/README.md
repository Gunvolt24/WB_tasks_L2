# unix-sort (L2.10)

Упрощённый аналог UNIX-утилиты `sort` c поддержкой флагов и внешней сортировки (разбиение на чанки и k-way merge).

---

## Структура проекта

```
┣ bin/                # бинарные файлы
┃ ┣ genbigtext
┃ ┗ sort
┣ cmd/
┃ ┣ genbigtext/       # генератор больших файлов
┃ ┗ main.go           # CLI sort
┣ e2e/                # интеграционные тесты CLI
┣ internal/           # пакеты: checker, chunk, comparator, inputoutput, merger, options, parser
┣ testdata/           # входные тестовые файлы (и big.txt)
┗ go.mod
```
---
## Поддерживаемые флаги

```bash
* `-k N` — сортировка по столбцу `N` (1..N). По умолчанию разделитель - TAB.
* `-d SEP` — разделитель столбцов (по умолчанию `\t`).
* `-n` — числовая сортировка.
* `-h` — сортировка  по числовому значению с учётом суффиксов (`512`, `1K`, `900M`, `3.5G` и т.п., основание 1024).
* `-M` — сортировка по названию месяца `Jan..Dec` (регистронезависимо).
* `-r` — сортировка в обратном порядке.
* `-u` — только уникальные строки (по ключу, зависящему от `-k`, `-d`, `-n/-h/-M`, `-b`).
* `-b` — игнорировать хвостовые пробелы у ключа.
* `-c` — проверить отсортированность потока (код выхода `0` — ок (без вывода в stdout), `1` — не отсортировано (вывод ошибки stderr)).
* `-chunk N` — максимум строк в чанке перед сбросом на диск (по умолчанию `1_000`).
* `-tmpdir DIR` — директория для временных файлов (по умолчанию системная TMP).
```

Проверка на наличие табуляции в файле:
```bash
`cat -A file.txt` (TAB отображается как `^I`).
```
---

## Сборка

Из корня проекта (`L2.10/`):

Bash:

```bash
go build -o bin/sort ./cmd
go build -o bin/genbigtext ./cmd/genbigtext
```

PowerShell:

```powershell
go build -o .\bin\sort.exe .\cmd
go build -o .\bin\genbigtext.exe .\cmd\genbigtext
```

---

## Генерация большого файла

Генератор создает файл с TAB-разделёнными столбцами: `id    month    human_size    number    key(with trailing spaces)`.

```bash
# 2000 строк в testdata/big.txt (пример)
./bin/genbigtext -n 2000 -o testdata/big.txt
```

Параметры:

* `-n` — количество строк (по умолчанию 1\_000),
* `-o` — путь к файлу (по умолчанию `testdata/big.txt`).

---

## Быстрый старт: 

### Примеры запуска с использованием файлов из `testdata/`

### Месяцы (столбец 2, TAB)

```bash
./bin/sort -k 2 -M testdata/months_sizes.txt
./bin/sort -k 2 -M -u testdata/months_sizes.txt          # уникально по месяцу
./bin/sort -k 2 -M -r testdata/months_sizes.txt          # обратный порядок
./bin/sort -k 2 -M testdata/months_sizes.txt | ./bin/sort -k 2 -M -c   # самопроверка
```

### "Человекочитаемые" размеры (столбец 3)

```bash
./bin/sort -k 3 -h testdata/months_sizes.txt             # возрастание
./bin/sort -k 3 -h -r testdata/months_sizes.txt          # убывание
./bin/sort -k 3 -h testdata/months_sizes.txt | ./bin/sort -k 3 -h -c
```

### Числовая сортировка (один столбец)

```bash
./bin/sort -n testdata/numbers.txt
./bin/sort -n -r testdata/numbers.txt
./bin/sort -n testdata/numbers.txt | ./bin/sort -n -c
```

### Уникальность по ключу

```bash
./bin/sort -k 2 -u testdata/dups_by_col2.txt             # по 2-й столбец
./bin/sort -k 3 -h -u testdata/dups_by_col2.txt          # по числовому значению с учетом суффиксов в 3-ем столбце
```

### Хвостовые пробелы в ключе

```bash
./bin/sort -k 2 -b -u testdata/trailing_blanks.txt       # останется одна строка "x<TAB>bar"
```

### Отсутствующие столбцы

```bash
./bin/sort -k 3 testdata/missing_cols.txt
./bin/sort -k 3 testdata/missing_cols.txt | ./bin/sort -k 3 -c
```

### Смешанные числовые/нечисловые значения

```bash
./bin/sort -k 2 -n testdata/mixed_non_numeric.txt
```

### Разделитель - пробел или запятая

```bash
./bin/sort -d " " -k 2 -M testdata/space_separated.txt
./bin/sort -d "," -k 2 -n testdata/data.csv
```

---

## Большие файлы и внешняя сортировка

Чтобы задействовать внешнюю сортировку, сделайте размер чанка меньше количества строк входа - тогда получится несколько временных файлов.

**Принцип работы:**
**1. Чтение и нарезка на чанки.**
- Берём вход `big.txt` (2000 строк).
- Копим, например, до 1000 строк в памяти `-chunk 1000` -> сортируем этот чанк по столбцу -`k N` и установленному флагу, например, `-M` -> пишем во временный файл в `./testdata/tmp/sort_chunk_*.txt`.
- Повторяем для следующей тысячи. В итоге получится N временных файла в зависимости от установленного разделения на чанки.

**2. Слияние (k-way merge)**
- Открываем все чанк-файлы, сливаем их в один отсортированный поток (по тем же правилам сортировки).
- Пишем результат либо в `stdout`, либо перенаправляем его в выходной файл.

**3. Отчистка**
- Временные файлы из `./testdata/tmp` по завершению сортировки удаляются, поэтому каталог остаётся пустым.

```bash
# По месяцам -> вывод в файл out_months.txt в корень проекта
./bin/sort -k 2 -M -chunk 1000 -tmpdir ./testdata/tmp ./testdata/big.txt > out_months.txt
# Проверка, что результат отсортирован по тем же правилам
./bin/sort -k 2 -M -c out_months.txt

# Аналогично для "человекочитаемых" числовых значений с суффиксом (сортировка по убыванию)
./bin/sort -k 3 -h -r -chunk 1000 -tmpdir ./testdata/tmp ./testdata/big.txt > out_human_desc.txt
./bin/sort -k 3 -h -c out_human_desc.txt

# Числовая сортировка по 4-му столбцу
./bin/sort -k 4 -n -chunk 1000 -tmpdir ./testdata/tmp ./testdata/big.txt > out_num.txt
./bin/sort -k 4 -n -c out_num.txt
```
---

## Чтение из STDIN

```bash
cat testdata/months_sizes.txt | ./bin/sort -k 2 -M
```

PowerShell:

```powershell
Get-Content .\testdata\months_sizes.txt | .\bin\sort.exe -k 2 -M
```

---

## Тесты и статический анализ

```bash
go test ./...
go vet ./...
golint ./...
```
