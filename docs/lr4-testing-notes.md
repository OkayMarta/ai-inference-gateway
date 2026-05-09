# LR4 Testing Notes

## Build/Test Results

### backend
Result: passed.

Command:

```bash
cd backend
go test ./...
```

Notes:
- Initial sandboxed run executed packages but exited with `go: failed to trim cache: open C:\Users\Asus\AppData\Local\go-build\trim.txt: Access is denied`.
- Re-run with normal Go cache access passed.
- Tested packages: handlers, repositories, services. Several packages have no test files.

### billing-service
Result: passed.

Command:

```bash
cd services/billing-service
go test ./...
```

Notes:
- Initial sandboxed run hit the same Go build cache permission issue.
- Re-run passed.
- Current service has no test files.

### task-service
Result: passed.

Command:

```bash
cd services/task-service
go test ./...
```

Notes:
- Initial sandboxed run hit the same Go build cache permission issue.
- Re-run passed.
- Current service has no test files.

### gateway-service
Result: passed.

Command:

```bash
cd services/gateway-service
go test ./...
```

Notes:
- Initial sandboxed run hit the same Go build cache permission issue.
- Re-run passed.
- Current service has no test files.

### frontend
Result: passed.

Command:

```bash
cd frontend
npm run build
```

Notes:
- Vite production build completed successfully.
- Frontend API base URL is `http://localhost:8080`.
- No direct frontend source calls to `http://localhost:8081` or `http://localhost:8082` were found.

## Configuration Checks

- `services/gateway-service/.env.example` exists and points to gateway `8080`, billing `8081`, task `8082`, and frontend `5173`.
- `services/billing-service/.env.example` exists and points to `billing_db`, port `8081`, `JWT_SECRET=dev-secret`, SMTP placeholders, frontend URL, and password reset TTL.
- `services/task-service/.env.example` exists and points to `task_db`, port `8082`, billing `8081`, and Ollama `11434`.
- Root `.gitignore` ignores `.env` and `.env.*`, while allowing `.env.example`.
- Real `.env` files exist locally in service folders and old `backend/`, but `git ls-files` shows only `.env.example` files are tracked.

## Static Inspection Summary

- Gateway exposes `/healthz`, public auth routes, public `/api/models`, and JWT-protected `/api/auth/me` plus task routes.
- Gateway validates HMAC JWT using `JWT_SECRET`, then forwards user identity to task-service via `X-User-ID`, `X-User-Email`, and `X-User-Role`.
- Gateway returns `503` with `billing service unavailable` or `task service unavailable` when the downstream HTTP client cannot connect.
- Billing-service owns register/login/forgot/reset, user lookup, charge, and refund logic.
- Password reset tokens are generated with random bytes, stored as SHA-256 hashes, checked for expiry and unused state, and marked used after reset.
- Task-service creates tasks by checking model existence, checking billing user existence, charging billing, then creating the task and refunding on task-create failure.
- Task-service starts a background worker, ensures default `worker-1`, refreshes `worker_supported_models`, and polls queued tasks every 500 ms.
- Frontend auth/task/model calls go through the gateway compatibility API client.

## Manual Test Checklist

- [ ] Start PostgreSQL and verify `billing_db` and `task_db` exist.
- [ ] Apply billing-service migrations `001_billing_init.sql` and `002_password_reset_tokens.sql`.
- [ ] Apply task-service migration `001_task_init.sql`.
- [ ] Start Ollama on `11434` and verify at least one model is installed.
- [ ] Start billing-service on `8081`.
- [ ] Start task-service on `8082`.
- [ ] Start gateway-service on `8080`.
- [ ] Start frontend on `5173`.
- [ ] `GET http://localhost:8080/healthz` returns `200 OK`.
- [ ] `GET http://localhost:8081/healthz` returns `200 OK`.
- [ ] `GET http://localhost:8082/healthz` returns `200 OK`.
- [ ] Register through gateway: `POST http://localhost:8080/api/auth/register`.
- [ ] Login through gateway: `POST http://localhost:8080/api/auth/login`.
- [ ] `GET http://localhost:8080/api/auth/me` with Bearer token returns the current user.
- [ ] Forgot password through gateway returns the generic reset message.
- [ ] Reset password through gateway with a valid token changes the password and consumes the token.
- [ ] `GET http://localhost:8080/api/models` returns active Ollama-synced models.
- [ ] Create task through gateway: `POST http://localhost:8080/api/tasks`.
- [ ] Confirm task appears in `task_db.prompt_tasks`.
- [ ] Confirm charge transaction appears in `billing_db.transactions`.
- [ ] Confirm user balance decreases in `billing_db.users`.
- [ ] Poll `GET http://localhost:8080/api/tasks` and observe `Queued -> Processing -> Completed` or `Failed`.
- [ ] Confirm completed task has result text from Ollama.
- [ ] Create a task and cancel while still `Queued`; confirm task becomes `Cancelled`.
- [ ] Confirm cancel refund transaction appears in `billing_db.transactions`.
- [ ] Stop task-service and verify gateway task routes return `503 task service unavailable`.
- [ ] Stop billing-service and verify gateway auth/me or task creation paths return a clear unavailable/downstream error.
- [ ] Exercise frontend register/login/dashboard/task submit/cancel/logout flow.
- [ ] Verify frontend browser network calls only target gateway `8080`.
- [ ] Verify `JWT_SECRET` values match in gateway-service and billing-service environments.
- [ ] Verify no real `.env` files are tracked before commit.

## Bugs Found

### Bug 1: internal billing endpoints can mutate balances without authentication if billing-service is reachable directly
Status: open
Area: billing-service / service boundary
Steps to reproduce:
1. Start billing-service on `8081`.
2. Send `POST http://localhost:8081/internal/billing/charge` or `/internal/billing/refund` with any existing `userId`, arbitrary `taskId`, and positive `amount`.
Expected:
Only task-service or trusted internal callers can charge/refund balances.
Actual:
The internal billing routes are mounted without auth or shared-service authentication.
Suggested fix:
Add service-to-service authentication for `/internal/*` routes, for example an internal shared secret header configured in both billing-service and task-service, or bind internal services to a private network only and enforce it in deployment.

### Bug 2: task-service trusts caller-supplied gateway identity headers if port 8082 is reachable directly
Status: open
Area: task-service / gateway boundary
Steps to reproduce:
1. Start task-service on `8082`.
2. Call task routes directly with `X-User-ID` and optionally `X-User-Role`.
Expected:
Only gateway can supply trusted identity headers to task-service.
Actual:
Task-service accepts any direct caller that provides `X-User-ID`.
Suggested fix:
Add service-to-service authentication between gateway and task-service, or ensure task-service is not exposed outside a private network. If keeping local ports public for the lab, document this as a local-only trust model.

### Bug 3: tasks can remain stuck in Queued when no active model mappings exist
Status: open
Area: task-service / worker startup / worker_supported_models
Steps to reproduce:
1. Start task-service with an empty `task_db` and Ollama unavailable, or with no active models.
2. Submit a task after models/mappings are missing or stale.
3. Observe worker polling.
Expected:
The service should clearly fail startup, reject task creation, or recover mappings so tasks are not silently stuck.
Actual:
Startup logs Ollama sync failures and continues. `RefreshSupportedModels` maps all workers to whatever active models exist. If no models are active or mappings are empty, `GetNextQueued` receives an empty supported model list and returns no task, so queued tasks may never be picked up.
Suggested fix:
On startup, after sync/mapping refresh, verify there is at least one active model and at least one worker/model mapping. If not, return a clear startup error or expose an unhealthy status. A smaller alternative is to reject task creation when no worker supports the selected model.

### Bug 4: cancelling a queued task refunds before persisting Cancelled status
Status: open
Area: task-service / billing compensation
Steps to reproduce:
1. Create a queued task.
2. Trigger cancellation while task DB update fails after billing refund succeeds.
Expected:
The system should not leave a refunded task still queued for execution.
Actual:
`CancelTask` calls billing refund first, then updates task status to `Cancelled`. If the DB update fails after refund, the code returns an error but the task may remain queued while the user has been refunded.
Suggested fix:
Use a safer compensation strategy. Options include marking the task cancelled first with a recoverable pending-refund state, adding an idempotent refund/compensation workflow, or re-charging/recording a reconciliation event if status persistence fails.

### Bug 5: direct public billing user lookup by id is unauthenticated
Status: open
Area: billing-service / user data access
Steps to reproduce:
1. Start billing-service on `8081`.
2. Call `GET http://localhost:8081/api/users/{id}` for a known user id.
Expected:
User lookup by id should require authentication or be internal-only.
Actual:
`/api/users/{id}` is mounted outside the authenticated group. Password hash is not serialized, but email, balance, role, and created time are exposed.
Suggested fix:
Move public `/api/users/{id}` into the authenticated group with access checks, or remove it and keep only `/internal/users/{id}` protected by service-to-service auth.

### Bug 6: whitespace-only task payload can be accepted through the API
Status: fixed
Area: task-service / validation
Steps to reproduce:
1. Submit a task through gateway with `"payload": "   "`.
Expected:
The API should reject empty or whitespace-only payloads.
Actual:
The handler checks only `req.Payload == ""`; `SubmitPrompt` later trims and stores an empty payload.
Suggested fix:
Trim payload before validation in `TaskHandler.Submit`, or validate the trimmed payload inside `InferenceService.SubmitPrompt`.
Fix applied:
`TaskHandler.Submit` now trims `modelId` and `payload` before validation and passes the trimmed values to `SubmitPrompt`.
