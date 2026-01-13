# HTTP-сервер «Календарь»

## Описание
Проект представляет собой простой сервис календаря событий, написанный на Go.
Реализован REST API для управления событиями календаря и минимальный frontend на HTML + JavaScript.
Проект построен с разделением на слои (domain, usecase, transport, repository) и придерживается принципов чистой архитектуры.

## Запуск проекта

**1.** Убедитесь, что установлен Go (версия 1.21 или выше)
**2.** Перейдите в корень проекта
**3.** Исправьте файл `.env.example` на `.env` и при необходимости измените параметры на:
```bash
HOST=localhost
PORT=8080
```
**4.** Запустите сервер командой:
```bash
go run cmd/main.go
```
После запуска:

- Frontend будет доступен по адресу: http://localhost:8080/
- REST API будет доступен по тем же адресам (например /create_event)

## Структура проекта
```text
┣ 📂cmd
┃ ┗ 📜main.go                       # main-сервис + статика
┣ 📂internal
┃ ┣ 📂domain
┃ ┃ ┣ 📜calendar.go                 # Доменная модель
┃ ┃ ┗ 📜errors.go                   # Доменные ошибки (валидация, не найдено и т.п.)
┃ ┣ 📂dto
┃ ┃ ┗ 📜dto.go                      # DTO для слоя usecase (входные данные)
┃ ┣ 📂transport
┃ ┃ ┗ 📂rest
┃ ┃   ┣ 📜error_mapper.go           # Маппинг доменных ошибок в HTTP-ответы
┃ ┃   ┣ 📜handler.go                # HTTP-хендлеры (обработка запросов)
┃ ┃   ┣ 📜request.go                # DTO для HTTP-запросов
┃ ┃   ┣ 📜response.go               # DTO для HTTP-ответов
┃ ┃   ┗ 📜router.go                 # Настройка маршрутов и middleware
┃ ┗ 📂usecase
┃   ┣ 📜calendar_service_test.go
┃   ┗ 📜calendar_service.go         # Бизнес-логика работы с календарём
┣ 📂middleware
┃ ┗ 📜logging.go                    # HTTP-middleware для логирования запросов  
┣ 📂repo
┃ ┣ 📜calendar_repo.go              # Интерфейс репозитория календаря
┃ ┣ 📜memory_test.go
┃ ┗ 📜memory.go                     # In-memory реализация репозитория
┣ 📂static
┃ ┣ 📜app.js                        # Клиентская логика (работа с API)
┃ ┗ 📜index.html                    # Простой фронтенд (страница календаря)
┣ 📜.env.example                    # Пример файла окружения
```