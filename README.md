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

<img width="759" height="115" alt="image" src="https://github.com/user-attachments/assets/fd860924-9bd8-4da2-a497-42114170a5e7" />


### 503 Service Unavailable — БД зупинена

```bash
docker compose stop postgres
curl -i localhost:8080/health
```

<img width="1324" height="552" alt="image" src="https://github.com/user-attachments/assets/7f42091d-bee5-4779-abf8-dcbaaa479b93" />

## Приклад логів

JSON-логи під час запуску застосунку:

```json
{"time":"2026-02-22T16:54:40.222852926Z","level":"INFO","msg":"migrations applied"}
{"time":"2026-02-22T16:54:40.225410801Z","level":"INFO","msg":"connected to database"}
{"time":"2026-02-22T16:54:40.225469676Z","level":"INFO","msg":"server started","addr":":8080"}
{"time":"2026-02-22T16:56:10.438844718Z","level":"INFO","msg":"request","method":"GET","path":"/health","status":200,"latency_ms":0,"client_ip":"192.168.65.1"}
{"time":"2026-02-22T16:56:23.512723918Z","level":"INFO","msg":"request","method":"GET","path":"/health","status":503,"latency_ms":53,"client_ip":"192.168.65.1"}
{"time":"2026-02-22T16:56:47.195966971Z","level":"INFO","msg":"request","method":"GET","path":"/health","status":200,"latency_ms":10,"client_ip":"192.168.65.1"}
{"time":"2026-02-22T17:00:59.133121046Z","level":"INFO","msg":"request","method":"GET","path":"/health","status":503,"latency_ms":5,"client_ip":"192.168.65.1"}
```
<img width="1431" height="417" alt="image" src="https://github.com/user-attachments/assets/9ec0f4dd-aff8-4c26-8114-81cdf17c09ae" />


## Підтвердження Graceful Shutdown

```bash
# Знайти PID процесу
lsof -i :8080

# Надіслати сигнал завершення
kill <pid>
```

Логи після надсилання сигналу:

```json
{"time":"2026-02-22T17:03:08.060873008Z","level":"INFO","msg":"shutting down server"}
{"time":"2026-02-22T17:03:08.060955342Z","level":"INFO","msg":"server exited properly"}
```

<img width="1437" height="432" alt="image" src="https://github.com/user-attachments/assets/a2396e63-0544-44e8-b073-57076464b0fa" />

