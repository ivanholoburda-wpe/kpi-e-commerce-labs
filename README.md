# Приклад сервісу для реєстрації та автентифікації з OTP/TOTP

REST API на Go для реєстрації користувачів з багатофакторною автентифікацією (OTP через Telegram, TOTP, PIN-код).
Цей репо буде використаний для лабораторних робіт з розгортання e-commerce та Системи Безпеки Програм і Даних.

## Налаштування

### Змінні оточення

| Змінна           | Обов'язкова | Опис                                          | Приклад                                                          |
|------------------|-------------|-----------------------------------------------|------------------------------------------------------------------|
| `DATABASE_URL`   | Так         | PostgreSQL connection string                   | `postgres://postgres:postgres@localhost:5432/lab2?sslmode=disable` |
| `TELEGRAM_TOKEN` | Так         | Токен Telegram-бота для відправки OTP          | `123456789:AAHdqTc...`                                           |
| `JWT_SECRET`     | Так         | Секретний ключ для підпису JWT-токенів         | `my-super-secret-key-256bit`                                     |
| `ADDR`           | Ні          | Адреса, на якій слухає сервер (за замовч. `:8080`) | `:8080`                                                          |

### Запуск

```bash
# 1. Запустити базу даних
docker compose up -d postgres

# 2. Встановити змінні оточення
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/lab2?sslmode=disable"
export TELEGRAM_TOKEN="<ваш токен>"
export JWT_SECRET="super-secret-key-change-me"

# 3. Запустити застосунок
go run ./cmd/main.go
```

Або повністю через Docker Compose:

```bash
TELEGRAM_TOKEN="<ваш токен>" docker compose up --build
```

## Підтвердження Health Check

### 200 OK — БД підключена

```
curl -i localhost:8080/health
```

<!-- TODO: скриншот терміналу з відповіддю 200 OK -->
![Health Check 200](screenshots/health-200.png)

### 503 Service Unavailable — БД зупинена

```bash
docker compose stop postgres
curl -i localhost:8080/health
```

<!-- TODO: скриншот терміналу з відповіддю 503 -->
![Health Check 503](screenshots/health-503.png)

## Приклад логів

JSON-логи під час запуску застосунку:

```json
{"time":"2026-02-16T12:00:00.000000+02:00","level":"INFO","msg":"migrations applied"}
{"time":"2026-02-16T12:00:00.001000+02:00","level":"INFO","msg":"connected to database"}
{"time":"2026-02-16T12:00:00.002000+02:00","level":"INFO","msg":"server started","addr":":8080"}
```

## Підтвердження Graceful Shutdown

```bash
# Знайти PID процесу
lsof -i :8080

# Надіслати сигнал завершення
kill <pid>
```

Логи після надсилання сигналу:

```json
{"time":"2026-02-16T12:05:00.000000+02:00","level":"INFO","msg":"shutting down server"}
{"time":"2026-02-16T12:05:00.005000+02:00","level":"INFO","msg":"server exited properly"}
```

<!-- TODO: скриншот терміналу з логами graceful shutdown -->
![Graceful Shutdown](screenshots/graceful-shutdown.png)
