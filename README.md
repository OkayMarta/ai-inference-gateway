# AI Inference Gateway (Lab 3)

## Опис

AI Inference Gateway — це навчальний full-stack проєкт, у якому backend на Go виступає шлюзом до локальних AI-моделей через Ollama, а frontend на React дає мінімальний інтерфейс для роботи із задачами.

Поточний стан проєкту відповідає **Lab 3**:

- backend працює з **PostgreSQL**
- моделі більше не є статичними або фейковими
- список моделей синхронізується з **Ollama**
- задачі обробляються **асинхронно** через background worker
- фінансові операції та створення/скасування задач виконуються **транзакційно**

## Основні можливості

- перегляд користувачів і їхнього токен-балансу
- перегляд доступних моделей із таблиці `ai_models`
- створення задач на генерацію через `POST /api/tasks`
- асинхронна обробка задач воркерами
- оновлення payload задачі, якщо вона ще в статусі `Queued`
- бізнес-скасування задачі через `DELETE /api/tasks/{id}`:
  - задача не видаляється фізично
  - статус змінюється на `Cancelled`
  - користувачу повертаються токени
  - створюється refund-транзакція
- фільтрація, пагінація та сортування задач

## Архітектура

Проєкт побудований як моноліт із шаровою структурою:

- `handlers` — HTTP-рівень
- `services` — бізнес-логіка
- `repositories` — доступ до PostgreSQL через `database/sql`
- `models` — доменні моделі

### Джерело істини для моделей

Архітектурне правило Lab 3:

- **Ollama** є джерелом істини для доступних моделей
- таблиця **`ai_models`** є persisted mirror/cache останньої успішної синхронізації
- `GET /api/models` читає **тільки** з PostgreSQL
- фейкові/default моделі більше **не використовуються**

### Поведінка при старті backend

Під час старту backend:

1. підключається до PostgreSQL
2. створює репозиторії
3. створює Ollama client
4. пробує виконати `SyncFromOllama()`
5. якщо sync успішний:
   - оновлює `ai_models`
   - перебудовує `worker_supported_models`
6. якщо sync неуспішний:
   - логуюється помилка
   - backend **все одно стартує**
   - останній валідний стан `ai_models` у БД зберігається

## Стек технологій

- **Backend:** Go, `chi`, `database/sql`, `lib/pq`
- **Frontend:** React, Vite
- **Database:** PostgreSQL
- **LLM runtime:** Ollama

## Структура проєкту

```text
backend/
  go.mod
  go.sum
  .env
  main.go
  migrations/
    001_init.sql
    002_seed.sql
  internal/
    db/
      postgres.go
    handlers/
    models/
    repositories/
    services/

frontend/
  src/
    api/
    components/
    styles/

scripts/
```

### Backend

#### `backend/main.go`

Composition root застосунку. Саме тут:

- завантажується `.env`
- ініціалізується PostgreSQL
- створюються репозиторії
- створюється Ollama client
- виконується startup sync моделей
- оновлюються worker-model mappings
- створюються сервіси та handlers
- запускається `WorkerService`
- конфігурується HTTP router

#### `backend/.env`

Локальна конфігурація середовища для backend:

- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `DB_SSLMODE`

Файл не комітиться в git.

#### `backend/go.mod`, `backend/go.sum`

Файли керування Go-залежностями backend-проєкту.

---

### Міграції

#### `backend/migrations/001_init.sql`

Створює схему бази даних:

- `users`
- `ai_models`
- `prompt_tasks`
- `transactions`
- `worker_nodes`
- `worker_supported_models`

Також містить:

- foreign keys
- check constraints
- індекси для задач, транзакцій і worker/model mapping

#### `backend/migrations/002_seed.sql`

Початкові seed-дані для Lab 3:

- тестові користувачі
- worker nodes

У цій міграції більше **немає** фейкових AI-моделей.

---

### `backend/internal/db`

#### `backend/internal/db/postgres.go`

Відповідає за підключення до PostgreSQL:

- читає конфігурацію з environment variables
- формує connection string
- відкриває `sql.DB`
- перевіряє підключення через `Ping()`

Також тут визначений `DBTX` — невеликий інтерфейс, який дозволяє репозиторіям працювати і з `*sql.DB`, і з `*sql.Tx`.

---

### `backend/internal/models`

#### `backend/internal/models/models.go`

Містить доменні моделі системи:

- `User`
- `AIModel`
- `PromptTask`
- `Transaction`
- `WorkerNode`

Тут же зберігаються доменні enum-like значення:

- статуси задач (`Queued`, `Processing`, `Completed`, `Failed`, `Cancelled`)
- статуси воркерів (`Idle`, `Busy`)

---

### `backend/internal/repositories`

Шар доступу до даних. Репозиторії працюють напряму з PostgreSQL через `database/sql` і не містять бізнес-логіки.

#### `backend/internal/repositories/common.go`

Спільні допоміжні функції для репозиторіїв, наприклад перевірка `RowsAffected`.

#### `backend/internal/repositories/user_repo.go`

Робота з таблицею `users`:

- читання користувачів
- читання по ID
- створення
- оновлення
- видалення
- оновлення балансу

#### `backend/internal/repositories/model_repo.go`

Робота з таблицею `ai_models`:

- отримання списку моделей
- отримання моделі по ID
- CRUD-операції
- `ReplaceAll(...)` для повної заміни каталогу моделей після sync з Ollama

#### `backend/internal/repositories/task_repo.go`

Робота з таблицею `prompt_tasks`:

- отримання задачі по ID
- список задач із фільтрами
- створення та оновлення задач
- `Complete(...)`, `Fail(...)`
- `GetNextQueued(...)` з PostgreSQL-safe вибіркою через `FOR UPDATE SKIP LOCKED`

#### `backend/internal/repositories/transaction_repo.go`

Робота з таблицею `transactions`:

- створення charge/refund транзакцій

#### `backend/internal/repositories/worker_repo.go`

Робота з:

- `worker_nodes`
- `worker_supported_models`

Функції:

- отримання воркерів
- отримання тільки idle воркерів
- оновлення статусу воркера
- зчитування `SupportedModels`
- transactional rebuild mapping між воркерами і моделями

---

### `backend/internal/services`

Шар бізнес-логіки та orchestration. Саме тут приймаються бізнес-рішення, а не в handler-ах чи репозиторіях.

#### `backend/internal/services/errors.go`

Єдине місце для доменних помилок, наприклад:

- `ErrUserNotFound`
- `ErrModelNotFound`
- `ErrTaskNotFound`
- `ErrInsufficientBalance`
- `ErrTaskCannotBeUpdated`
- `ErrTaskCannotBeDeleted`
- `ErrInvalidPagination`

#### `backend/internal/services/repository_interfaces.go`

Інтерфейси репозиторіїв, від яких залежить service layer:

- `UserRepository`
- `ModelRepository`
- `TaskRepository`
- `TransactionRepository`
- `WorkerRepository`

Також тут знаходиться `TaskListFilter`.

#### `backend/internal/services/user_service.go`

Бізнес-логіка для користувачів:

- отримання списку користувачів
- отримання користувача по ID
- оновлення `username` / `tokenBalance`

#### `backend/internal/services/model_service.go`

Бізнес-логіка для моделей:

- читання моделей із БД
- `SyncFromOllama()`
- мапінг Ollama metadata у внутрішню `AIModel`

#### `backend/internal/services/ollama_client.go`

HTTP client для інтеграції з Ollama:

- `ListModels()`
- `Generate(...)`

#### `backend/internal/services/inference_service.go`

Основна бізнес-логіка задач:

- `SubmitPrompt(...)`
- `GetTaskByID(...)`
- `ListTasks(...)`
- `UpdateTaskPayload(...)`
- `CancelTask(...)`

Тут реалізовані транзакційні сценарії:

- створення задачі зі списанням токенів
- скасування задачі з поверненням токенів

#### `backend/internal/services/worker_service.go`

Фоновий worker orchestration:

- polling idle воркерів
- вибір сумісної queued задачі
- запуск реальної генерації через Ollama
- оновлення статусів задач і воркерів
- refresh worker/model mappings після startup sync моделей

---

### `backend/internal/handlers`

HTTP-рівень застосунку. Handler-и:

- читають параметри запиту
- валідовують body/query на базовому рівні
- викликають service layer
- формують JSON response

#### `backend/internal/handlers/response.go`

Спільні helper-и для HTTP-відповідей:

- `respondJSON(...)`
- `respondError(...)`
- `ErrorResponse`
- `mapErrorToStatus(...)`

#### `backend/internal/handlers/user_handler.go`

Endpoints для users:

- `GET /api/users`
- `GET /api/users/{id}`
- `PUT /api/users/{id}`

#### `backend/internal/handlers/model_handler.go`

Endpoint:

- `GET /api/models`

Читає моделі тільки з БД через service layer.

#### `backend/internal/handlers/task_handler.go`

Endpoints для задач:

- `POST /api/tasks`
- `GET /api/tasks`
- `GET /api/tasks/{id}`
- `PUT /api/tasks/{id}`
- `DELETE /api/tasks/{id}`

#### `backend/internal/handlers/recovery.go`

Recovery middleware для перехоплення panic і повернення коректної JSON-помилки.

---

### Frontend

Frontend залишається мінімальним і демонстраційним.

#### `frontend/src/api`

Код для HTTP-запитів до backend.

##### `frontend/src/api/client.js`

Містить функції:

- `getUsers()`
- `getModels()`
- `getTasks(...)`
- `submitTask(...)`
- `deleteTask(...)`

#### `frontend/src/components`

React-компоненти інтерфейсу.

##### `frontend/src/components/Dashboard.jsx`

Головний контейнер сторінки:

- завантажує users/models/tasks
- тримає state вибору user/model
- запускає polling задач
- координує submit і cancel

##### `frontend/src/components/TaskComposer.jsx`

Форма створення задачі:

- вибір user
- вибір model
- textarea для prompt
- submit

Також коректно обробляє стан, коли доступних моделей немає.

##### `frontend/src/components/TaskList.jsx`

Список задач:

- рендер карток задач
- status filter
- cancel для queued задач
- summary по статусах

##### Інші UI-компоненти

- `SectionCard.jsx` — візуальна оболонка секцій
- `EmptyState.jsx` — порожні стани
- `StatusBadge.jsx` — бейджі статусів задач

#### `frontend/src/styles`

CSS-стилі компонентів.

##### `frontend/src/styles/components/Dashboard.css`

Основні стилі dashboard:

- layout
- composer
- task list
- notices
- status/filter controls

---

### `scripts`

Допоміжна директорія для додаткових скриптів, якщо вони потрібні для локальної розробки або демонстрації.

## Підготовка середовища

### 1. PostgreSQL

Створи БД, наприклад:

```sql
CREATE DATABASE ai_inference_gateway;
```

Потім застосуй міграції:

- `backend/migrations/001_init.sql`
- `backend/migrations/002_seed.sql`

`002_seed.sql` зараз:

- додає тестових користувачів
- додає worker nodes
- **не** додає жодних fake AI models

### 2. `.env`

У директорії `backend/` має бути файл `.env`, наприклад:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ai_inference_gateway
DB_SSLMODE=disable
```

Backend читає `.env` через `godotenv.Load()`.

### 3. Ollama

Для реальної генерації потрібна запущена Ollama і хоча б одна завантажена модель.

Приклади команд:

```bash
ollama serve
ollama pull llama3:8b
ollama list
```

## Запуск проєкту

### Backend

```bash
cd backend
go run main.go
```

Backend стартує на:

```text
http://localhost:8080
```

Health endpoint:

```text
GET /healthz
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend стартує на:

```text
http://localhost:5173
```

## Актуальні API endpoints

### Users

| Method | Endpoint           | Description |
| ------ | ------------------ | ----------- |
| `GET`  | `/api/users`       | Отримати список користувачів |
| `GET`  | `/api/users/{id}`  | Отримати користувача за ID |
| `PUT`  | `/api/users/{id}`  | Оновити `username` або `tokenBalance` |

### Models

| Method | Endpoint      | Description |
| ------ | ------------- | ----------- |
| `GET`  | `/api/models` | Отримати моделі з PostgreSQL (`ai_models`) |

### Tasks

| Method   | Endpoint           | Description |
| -------- | ------------------ | ----------- |
| `POST`   | `/api/tasks`       | Створити нову задачу |
| `GET`    | `/api/tasks`       | Отримати список задач |
| `GET`    | `/api/tasks/{id}`  | Отримати задачу за ID |
| `PUT`    | `/api/tasks/{id}`  | Оновити payload, якщо задача ще `Queued` |
| `DELETE` | `/api/tasks/{id}`  | Скасувати задачу бізнес-способом |

### Доступні query params для `GET /api/tasks`

- `userId`
- `status`
- `limit`
- `offset`
- `sort`

Підтримувані значення `sort`:

- `created_at_desc`
- `created_at_asc`

Приклад:

```bash
curl "http://localhost:8080/api/tasks?userId=user-1&status=Completed&limit=10&offset=0&sort=created_at_desc"
```

## Приклади запитів

### Створення задачі

```bash
curl -X POST http://localhost:8080/api/tasks ^
  -H "Content-Type: application/json" ^
  -d "{\"userId\":\"user-1\",\"modelId\":\"llama3-8b\",\"payload\":\"Explain worker queues\"}"
```

### Оновлення користувача

```bash
curl -X PUT http://localhost:8080/api/users/user-1 ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"alice-updated\",\"tokenBalance\":120.5}"
```

### Оновлення payload задачі

```bash
curl -X PUT http://localhost:8080/api/tasks/task-123 ^
  -H "Content-Type: application/json" ^
  -d "{\"payload\":\"new prompt text\"}"
```

### Скасування задачі

```bash
curl -X DELETE http://localhost:8080/api/tasks/task-123
```

## Життєвий цикл задачі

Підтримувані статуси:

- `Queued`
- `Processing`
- `Completed`
- `Failed`
- `Cancelled`

Типовий flow:

1. користувач створює задачу
2. backend транзакційно:
   - перевіряє користувача
   - перевіряє модель
   - списує токени
   - створює задачу
   - створює charge-транзакцію
3. background worker бере найстарішу сумісну задачу
4. задача переходить у `Processing`
5. якщо Ollama повертає реальний результат:
   - задача переходить у `Completed`
6. якщо Ollama недоступна або generation падає:
   - задача переходить у `Failed`
7. якщо користувач скасовує задачу в `Queued`:
   - задача переходить у `Cancelled`
   - токени повертаються
   - створюється refund-транзакція

## Data Flow

Нижче — основні наскрізні сценарії в системі. Вони добре показують, як взаємодіють шари `handler -> service -> repository -> DB/Ollama`.

### 1. Отримання списку моделей

Flow для `GET /api/models`:

1. `ModelHandler.GetAll()` приймає HTTP-запит
2. викликає `ModelService.GetAll()`
3. `ModelService` делегує читання в `ModelRepository.GetAll()`
4. `ModelRepository` виконує SQL `SELECT` з таблиці `ai_models`
5. результат повертається назад у handler
6. handler віддає JSON-відповідь клієнту

Важливо:

- під час `GET /api/models` **немає прямого виклику Ollama**
- API читає тільки persisted каталог із PostgreSQL

### 2. Startup sync моделей з Ollama

Flow під час старту backend:

1. `main.go` створює `ModelService` і `WorkerService`
2. `main.go` викликає `modelSvc.SyncFromOllama()`
3. `ModelService` звертається до `OllamaClient.ListModels()`
4. отримані Ollama-моделі мапляться у внутрішні `AIModel`
5. `ModelRepository.ReplaceAll(...)` транзакційно оновлює `ai_models`
6. якщо sync успішний, `main.go` викликає `workerSvc.RefreshSupportedModels()`
7. `WorkerService` читає актуальні моделі з БД
8. `WorkerRepository.ReplaceSupportedModelsForAllWorkers(...)` перебудовує `worker_supported_models`

Важливо:

- якщо Ollama недоступна, startup sync просто логуює помилку
- backend **не падає**
- попередній валідний стан `ai_models` у БД зберігається

### 3. Створення задачі

Flow для `POST /api/tasks`:

1. `TaskHandler.Submit()` читає JSON body
2. handler викликає `InferenceService.SubmitPrompt(...)`
3. service відкриває DB transaction
4. через репозиторії всередині transaction:
   - читає користувача
   - читає модель
   - перевіряє баланс
   - списує токени
   - створює запис у `prompt_tasks`
   - створює charge-запис у `transactions`
5. transaction комітиться
6. handler повертає створену задачу клієнту

Важливо:

- створення задачі є **атомарним**
- немає часткового стану типу “баланс списався, але задача не створилась”

### 4. Асинхронна обробка задачі воркером

Flow background processing:

1. `WorkerService.Start()` запускає polling loop
2. `WorkerService.processNext()` отримує idle воркерів через `WorkerRepository.GetIdle()`
3. для кожного воркера викликається `TaskRepository.GetNextQueued(...)`
4. repository:
   - транзакційно знаходить найстарішу сумісну queued-задачу
   - блокує її через `FOR UPDATE SKIP LOCKED`
   - переводить у `Processing`
5. `WorkerService` ставить воркера в `Busy`
6. `executeTask(...)`:
   - читає модель з БД
   - викликає реальну генерацію через `OllamaClient.Generate(...)`
7. якщо Ollama повернула результат:
   - `TaskRepository.Complete(...)`
   - статус стає `Completed`
8. якщо Ollama недоступна або generation падає:
   - `TaskRepository.Fail(...)`
   - статус стає `Failed`
9. воркер повертається у `Idle`

Важливо:

- успішних fake/simulated результатів більше немає
- completed-задача завжди означає реальний результат від Ollama

### 5. Скасування задачі

Flow для `DELETE /api/tasks/{id}`:

1. `TaskHandler.DeleteTask()` отримує task ID
2. handler викликає `InferenceService.CancelTask(id)`
3. service відкриває DB transaction
4. всередині transaction:
   - читає задачу
   - перевіряє, що вона ще `Queued`
   - читає модель
   - читає користувача
   - повертає токени на баланс
   - змінює статус задачі на `Cancelled`
   - створює refund-транзакцію
5. transaction комітиться
6. handler повертає оновлену задачу клієнту

Важливо:

- задача не видаляється фізично
- історія зберігається в БД
- refund, зміна статусу і запис транзакції відбуваються **атомарно**

## Важливі правила поточної реалізації

### Немає фейкових моделей

- статичні/default моделі видалені
- якщо Ollama недоступна, backend не створює fake data
- `GET /api/models` повертає те, що є в `ai_models`
- якщо `ai_models` порожня, API повертає `[]`

### Немає fake generation

- WorkerService виконує тільки реальну генерацію через Ollama
- якщо Ollama недоступна, задача завершується зі статусом `Failed`
- псевдоуспішні симуляції більше не використовуються

### Тестова симуляція помилки

Для демонстрації failure flow можна використати marker:

```text
__SIMULATE_FAILURE__
```

Якщо `payload` містить цей marker, worker навмисно завершить задачу зі статусом `Failed`.

## Фронтенд

Frontend навмисно залишений простим.

Поточна UI-підтримка:

- вибір користувача
- вибір моделі
- створення задачі
- перегляд списку задач
- cancel для `Queued` задач
- фільтр задач за статусом

Якщо `/api/models` повертає порожній список:

- frontend показує notice про відсутність моделей
- submit форми блокується
- це вважається валідним станом системи

## Перевірка збірки

### Backend

```bash
cd backend
go test ./...
```

### Frontend

```bash
cd frontend
npm run build
```

## Короткий підсумок Lab 3

У Lab 3 проєкт перейшов від in-memory підходу до persisted архітектури:

- PostgreSQL замість in-memory storage
- SQL міграції
- транзакційні бізнес-операції
- синхронізація реальних моделей з Ollama
- persisted mapping між worker-ами і моделями
- реальна асинхронна обробка задач без fake success fallback
