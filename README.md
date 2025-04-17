# Bank API

RESTful API для управления банковскими счетами, картами, кредитами и аналитикой. Полностью покрывает требования проекта: регистрация, аутентификация, выпуск карт, переводы, кредиты, шедулер, SMTP и SOAP-интеграции, защита данных, аналитика.

## Запуск

### Зависимости
- Go 1.22+
- PostgreSQL 13+

### Конфигурация окружения `.env`
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=bank
JWT_SECRET=supersecretkey
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=user@example.com
SMTP_PASSWORD=password
```

### Команды
```
go mod tidy
make run
```

Приложение запускается на `http://localhost:8080`

## Аутентификация

JWT передается через заголовок:
```
Authorization: Bearer <token>
```

## Эндпоинты

### Аутентификация
- `POST /register` — регистрация
- `POST /login` — вход, возвращает JWT

### Счета и переводы
- `POST /accounts` — создание счета
- `POST /transfer` — перевод между счетами

### Карты
- `POST /cards` — выпуск виртуальной карты
- `GET /cards/{id}` — расшифровка карты

### Кредиты
- `POST /credits` — оформление кредита (аннуитет)
- `GET /credits/{creditId}/schedule` — график платежей

### Аналитика
- `GET /analytics` — агрегированные показатели
- `GET /accounts/{accountId}/predict?days=N` — прогноз баланса

## Шедулер
- Запускается каждые 12 часов
- Обрабатывает просроченные платежи, начисляет 10% штраф

## Интеграции
- SMTP: отправка уведомлений по e-mail
- SOAP: получение ключевой ставки из ЦБ РФ

## Безопасность
- JWT + Middleware
- bcrypt (пароли и CVV)
- OpenPGP + HMAC для шифрования номера карты и срока действия
- Контроль доступа на уровне пользователя

## Тестирование
```
make test       # Unit + интеграционные
make cover      # Покрытие кода (HTML)
```

Покрытие включает: сервисы, обработчики, middleware, репозитории, утилиты, шедулер

## Служебные файлы
- `Makefile` — запуск, тестирование, покрытие
- `docs/postman_collection.json` — коллекция для Postman
- `migrations/` — SQL-миграции
- `utils/pgp_utils.go` — реализация PGP + HMAC


