# Banking REST API

## Описание проекта

Banking REST API — это серверное приложение, реализующее базовые функции банковской системы: управление пользователями, банковскими счетами, виртуальными картами, переводами, кредитами, а также аналитикой. Реализованы интеграции с внешними сервисами (ЦБ РФ и SMTP), обеспечена безопасность на уровне аутентификации и защиты данных.

Проект написан на Go с использованием PostgreSQL в качестве базы данных.

## Возможности

- Регистрация и аутентификация пользователей (JWT)
- Создание банковских счетов и управление ими
- Пополнение баланса и переводы между счетами
- Генерация и просмотр виртуальных карт
- Шифрование номера карты (PGP), хеширование CVV (bcrypt), HMAC
- Кредитование: оформление, график платежей, начисление штрафов
- Аналитика: расходы, доходы, кредитная нагрузка, прогноз баланса
- Интеграция с API ЦБ РФ для получения ключевой ставки (SOAP)
- Уведомления по email (SMTP)
- Middleware для защиты маршрутов

## Используемые технологии

- Язык: Go 1.23+
- База данных: PostgreSQL 17 + pgcrypto
- Фреймворк маршрутизации: gorilla/mux
- JWT: github.com/golang-jwt/jwt/v5
- PGP: github.com/ProtonMail/go-crypto/openpgp
- SMTP: github.com/go-mail/mail/v2
- XML/SOAP: github.com/beevik/etree
- Хеширование: golang.org/x/crypto/bcrypt
- Логирование: github.com/sirupsen/logrus

## Структура проекта

banking-api/
├── cmd/server               # Точка входа и HTTP-обработчики
│   ├── handlers             # Контроллеры
│   └── middleware           # JWT Middleware
├── internal/db              # Репозитории и бизнес-логика
│   └── services             # SMTP, SOAP, Scheduler
├── models                   # Определения структур моделей
├── config                   # Конфигурации
├── db/migrations            # SQL-миграции
└── main.go

## Как запустить

1. Установите PostgreSQL и создайте базу данных:
   CREATE DATABASE banking;

2. Примените миграции:
   psql -U your_user -d banking -f db/migrations/1_create_tables.up.sql

3. Создайте PGP-пару ключей и поместите public_pgp_key.asc в ~/

4. Установите переменные окружения:
   export JWT_SECRET="your_secret"
   export SMTP_USER="your_email@example.com"
   export SMTP_PASS="your_password"

5. Запустите сервер:
   go run cmd/server/main.go

## Примеры запросов

### Регистрация
POST /register
{
  "email": "user@example.com",
  "password": "123456"
}

### Аутентификация
POST /login
{
  "email": "user@example.com",
  "password": "123456"
}

### Перевод между счетами
POST /transfer?from_account_id=1&to_account_id=2
Authorization: Bearer <JWT>
{
  "amount": 500.0
}

## Тестирование

Приложение можно тестировать с помощью curl, Postman или любого HTTP-клиента. JWT-токен необходимо передавать в заголовке Authorization.

## Автор

Dmitry Rubanov  

