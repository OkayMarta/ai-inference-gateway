# AI Inference Gateway (Lab 2)

## Опис

AI Inference Gateway — це навчальний full-stack проєкт для ЛР2, який імітує шлюз до AI-моделей.

Система складається з:

- backend на Go
- frontend на React
- in-memory storage
- фонового worker-а для обробки задач

У межах ЛР2 застосовано монолітну архітектуру: HTTP API, бізнес-логіка, черга задач і in-memory репозиторії працюють в одному backend-застосунку.

## Основний функціонал

- створення задач на обробку prompt для AI-моделі
- перевірка та списання балансу користувача
- фонова обробка задач зі статусами `Queued -> Processing -> Completed/Failed`
- перегляд списку задач
- перегляд користувачів
- перегляд доступних моделей

## Стек технологій

- Backend: Go, `net/http`, `chi`
- Frontend: React, Vite
- Storage: in-memory (`map` + `sync.RWMutex`)
- Обробка задач: background worker queue

## Запуск проєкту

### Backend

1. Перейдіть у директорію backend:

```bash
cd backend
```

2. Запустіть сервер:

```bash
go run main.go
```

3. Backend буде доступний за адресою:

```text
http://localhost:8080
```

Додатково:

- healthcheck: `GET http://localhost:8080/healthz`
- CORS налаштований для frontend на `http://localhost:5173`

### Frontend

1. Перейдіть у директорію frontend:

```bash
cd frontend
```

2. Встановіть залежності:

```bash
npm install
```

3. Запустіть dev-сервер:

```bash
npm run dev
```

4. Frontend буде доступний за адресою:

```text
http://localhost:5173
```

Для production build:

```bash
npm run build
```

## API

| Method | Endpoint                     | Description                             |
| ------ | ---------------------------- | --------------------------------------- |
| `GET`  | `/healthz`                   | Перевірка доступності backend           |
| `GET`  | `/api/users`                 | Отримати список користувачів            |
| `GET`  | `/api/users/{id}`            | Отримати користувача за ID              |
| `GET`  | `/api/models`                | Отримати список доступних моделей       |
| `POST` | `/api/tasks`                 | Створити нову задачу                    |
| `GET`  | `/api/tasks`                 | Отримати всі задачі                     |
| `GET`  | `/api/tasks?userId={userId}` | Отримати задачі конкретного користувача |
| `GET`  | `/api/tasks/{id}`            | Отримати задачу за ID                   |

## Приклади запитів

### 1. POST /api/tasks

Приклад body:

```json
{
    "userId": "user-1",
    "modelId": "model-1",
    "payload": "Explain how a worker queue works"
}
```

Поля:

- `userId` — ID користувача
- `modelId` — ID моделі
- `payload` — текст запиту для обробки

Приклад через `curl`:

```bash
curl -X POST http://localhost:8080/api/tasks ^
  -H "Content-Type: application/json" ^
  -d "{\"userId\":\"user-1\",\"modelId\":\"model-1\",\"payload\":\"Explain how a worker queue works\"}"
```

### 2. GET /api/tasks

Отримати всі задачі:

```bash
curl http://localhost:8080/api/tasks
```

Отримати задачі конкретного користувача:

```bash
curl "http://localhost:8080/api/tasks?userId=user-1"
```

Поведінка:

- без `userId` повертаються всі задачі
- з `userId` повертаються тільки задачі цього користувача

## Flow задачі

1. Користувач створює задачу через `POST /api/tasks`.
2. Backend перевіряє існування користувача і моделі.
3. Перевіряється баланс користувача і виконується списання.
4. Створюється задача зі статусом `Queued`.
5. Фоновий worker бере задачу в обробку і переводить її в `Processing`.
6. Після завершення задача отримує статус `Completed` або `Failed`.

## Симуляція помилки

Для тестування можна штучно викликати помилку обробки.

Якщо `payload` містить маркер:

```text
__SIMULATE_FAILURE__
```

задача буде завершена зі статусом `Failed`.

Це використовується тільки для тестування та демонстрації.

## Структура проєкту

```text
backend/
  internal/
    handlers/
    services/
    repositories/
    models/

frontend/
  src/
    api/
    components/
    styles/

scripts/
```

## Додаткові примітки

- backend використовує in-memory сховище, тому дані не зберігаються між перезапусками
- worker-и створюються під час старту backend
- якщо Ollama недоступна, backend використовує симуляцію моделей для демонстрації роботи системи
