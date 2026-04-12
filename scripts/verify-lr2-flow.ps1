param(
    [string]$BaseUrl = "http://127.0.0.1:8080",
    [int]$PollAttempts = 30,
    [int]$PollDelaySeconds = 1
)

$ErrorActionPreference = "Stop"

function Invoke-Api {
    param(
        [string]$Method,
        [string]$Path,
        [object]$Body = $null
    )

    $uri = "{0}{1}" -f $BaseUrl.TrimEnd("/"), $Path
    $headers = @{ "Content-Type" = "application/json" }
    $jsonBody = $null

    if ($null -ne $Body) {
        $jsonBody = $Body | ConvertTo-Json -Compress
    }

    try {
        if ($null -ne $jsonBody) {
            $response = Invoke-WebRequest -UseBasicParsing -Method $Method -Uri $uri -Headers $headers -Body $jsonBody
        } else {
            $response = Invoke-WebRequest -UseBasicParsing -Method $Method -Uri $uri -Headers $headers
        }

        return [pscustomobject]@{
            StatusCode = [int]$response.StatusCode
            RawBody = $response.Content
            Json = Convert-ToJsonIfPossible -RawBody $response.Content
        }
    } catch {
        $statusCode = [int]$_.Exception.Response.StatusCode
        $rawBody = $_.ErrorDetails.Message

        return [pscustomobject]@{
            StatusCode = $statusCode
            RawBody = $rawBody
            Json = Convert-ToJsonIfPossible -RawBody $rawBody
        }
    }
}

function Convert-ToJsonIfPossible {
    param([string]$RawBody)

    if (-not $RawBody) {
        return $null
    }

    $trimmed = $RawBody.Trim()
    if (-not ($trimmed.StartsWith("{") -or $trimmed.StartsWith("["))) {
        return $null
    }

    try {
        return $trimmed | ConvertFrom-Json
    } catch {
        return $null
    }
}

function Assert-True {
    param(
        [object]$Condition,
        [string]$Message
    )

    if (-not $Condition) {
        throw $Message
    }
}

function Assert-Status {
    param(
        [object]$Response,
        [int]$ExpectedStatus,
        [string]$Context
    )

    Assert-True ($Response.StatusCode -eq $ExpectedStatus) "$Context returned $($Response.StatusCode), expected $ExpectedStatus"
}

function Assert-ErrorFormat {
    param(
        [object]$Response,
        [int]$ExpectedStatus,
        [string]$Context
    )

    Assert-Status $Response $ExpectedStatus $Context
    Assert-True ($null -ne $Response.Json) "$Context did not return JSON body"
    Assert-True ($Response.Json.error) "$Context did not include error field"
    Assert-True ($Response.Json.message) "$Context did not include message field"
    Assert-True ($Response.Json.status -eq $ExpectedStatus) "$Context returned mismatched error status in JSON body"
}

function Wait-ForTaskStatus {
    param(
        [string]$TaskId,
        [string]$ExpectedStatus,
        [string]$Context
    )

    for ($attempt = 1; $attempt -le $PollAttempts; $attempt++) {
        $response = Invoke-Api -Method "GET" -Path "/api/tasks/$TaskId"
        Assert-Status $response 200 "$Context lookup"

        if ($response.Json.status -eq $ExpectedStatus) {
            return $response.Json
        }

        Start-Sleep -Seconds $PollDelaySeconds
    }

    throw "$Context did not reach status $ExpectedStatus after $PollAttempts attempts"
}

$report = New-Object System.Collections.Generic.List[string]

$health = Invoke-Api -Method "GET" -Path "/healthz"
Assert-Status $health 200 "GET /healthz"
Assert-True ($health.RawBody -eq "OK") "GET /healthz returned unexpected body"
$report.Add("healthz ok")

$usersResponse = Invoke-Api -Method "GET" -Path "/api/users"
Assert-Status $usersResponse 200 "GET /api/users"
$users = @($usersResponse.Json)
Assert-True ($users.Count -ge 3) "GET /api/users returned too few seeded users"
$report.Add("users ok")

$userLookup = Invoke-Api -Method "GET" -Path "/api/users/user-1"
Assert-Status $userLookup 200 "GET /api/users/user-1"
Assert-True ($userLookup.Json.id -eq "user-1") "GET /api/users/user-1 returned wrong user"
$report.Add("user lookup ok")

$modelsResponse = Invoke-Api -Method "GET" -Path "/api/models"
Assert-Status $modelsResponse 200 "GET /api/models"
$models = @($modelsResponse.Json)
Assert-True ($models.Count -ge 1) "GET /api/models returned no models"
$model = $models[0]
$report.Add("models ok")

$beforeSubmitBalance = ($users | Where-Object { $_.id -eq "user-1" }).tokenBalance
$happyPayload = "LR2 verification happy path $(Get-Date -Format o)"
$happySubmit = Invoke-Api -Method "POST" -Path "/api/tasks" -Body @{
    userId = "user-1"
    modelId = $model.id
    payload = $happyPayload
}
Assert-Status $happySubmit 201 "POST /api/tasks happy path"
Assert-True ($happySubmit.Json.status -eq "Queued") "Happy path task was not queued after submit"
$happyTaskId = $happySubmit.Json.id
$report.Add("task submit ok")

$queuedTask = Invoke-Api -Method "GET" -Path "/api/tasks/$happyTaskId"
Assert-Status $queuedTask 200 "GET /api/tasks/{id} for happy path"
Assert-True ($queuedTask.Json.id -eq $happyTaskId) "Task lookup returned unexpected id"

$completedTask = Wait-ForTaskStatus -TaskId $happyTaskId -ExpectedStatus "Completed" -Context "Happy path task"
Assert-True ($completedTask.result) "Completed task did not include result"
$report.Add("worker completed happy path")

$tasksByUser = Invoke-Api -Method "GET" -Path "/api/tasks?userId=user-1"
Assert-Status $tasksByUser 200 "GET /api/tasks?userId=user-1"
$matchingUserTasks = @(@($tasksByUser.Json) | Where-Object { $_.id -eq $happyTaskId })
Assert-True ($matchingUserTasks.Count -eq 1) "User task list did not include happy path task"
$report.Add("user task list ok")

$updatedUsers = Invoke-Api -Method "GET" -Path "/api/users"
Assert-Status $updatedUsers 200 "GET /api/users after submit"
$afterSubmitBalance = (@($updatedUsers.Json) | Where-Object { $_.id -eq "user-1" }).tokenBalance
$expectedBalance = [double]$beforeSubmitBalance - [double]$model.tokenCost
Assert-True ([double]$afterSubmitBalance -eq $expectedBalance) "User balance did not decrease by model cost after submit"
$report.Add("balance update ok")

$missingTask = Invoke-Api -Method "GET" -Path "/api/tasks/missing-task"
Assert-ErrorFormat $missingTask 404 "GET /api/tasks/missing-task"

$missingUserLookup = Invoke-Api -Method "GET" -Path "/api/users/missing-user"
Assert-ErrorFormat $missingUserLookup 404 "GET /api/users/missing-user"
$report.Add("not found responses ok")

$missingUserSubmit = Invoke-Api -Method "POST" -Path "/api/tasks" -Body @{
    userId = "missing-user"
    modelId = $model.id
    payload = "missing user"
}
Assert-ErrorFormat $missingUserSubmit 404 "POST /api/tasks with missing user"

$missingModelSubmit = Invoke-Api -Method "POST" -Path "/api/tasks" -Body @{
    userId = "user-1"
    modelId = "missing-model"
    payload = "missing model"
}
Assert-ErrorFormat $missingModelSubmit 404 "POST /api/tasks with missing model"

$missingFieldsSubmit = Invoke-Api -Method "POST" -Path "/api/tasks" -Body @{
    userId = "user-1"
}
Assert-ErrorFormat $missingFieldsSubmit 400 "POST /api/tasks with missing fields"

$emptyBodySubmit = Invoke-Api -Method "POST" -Path "/api/tasks" -Body @{}
Assert-ErrorFormat $emptyBodySubmit 400 "POST /api/tasks with empty body"
$report.Add("bad request responses ok")

$userTwo = Invoke-Api -Method "GET" -Path "/api/users/user-2"
Assert-Status $userTwo 200 "GET /api/users/user-2"

while ([double]$userTwo.Json.tokenBalance -ge [double]$model.tokenCost) {
    $prepSubmit = Invoke-Api -Method "POST" -Path "/api/tasks" -Body @{
        userId = "user-2"
        modelId = $model.id
        payload = "LR2 insufficient balance prep $(Get-Date -Format o)"
    }

    Assert-Status $prepSubmit 201 "POST /api/tasks insufficient balance prep"
    $userTwo = Invoke-Api -Method "GET" -Path "/api/users/user-2"
    Assert-Status $userTwo 200 "GET /api/users/user-2 after prep"
}

$insufficientBalanceSubmit = Invoke-Api -Method "POST" -Path "/api/tasks" -Body @{
    userId = "user-2"
    modelId = $model.id
    payload = "insufficient balance"
}
Assert-ErrorFormat $insufficientBalanceSubmit 422 "POST /api/tasks with insufficient balance"
$report.Add("unprocessable entity response ok")

$failedSubmit = Invoke-Api -Method "POST" -Path "/api/tasks" -Body @{
    userId = "user-3"
    modelId = $model.id
    payload = "[fail] LR2 failure path $(Get-Date -Format o)"
}
Assert-Status $failedSubmit 201 "POST /api/tasks failed path"

$failedTask = Wait-ForTaskStatus -TaskId $failedSubmit.Json.id -ExpectedStatus "Failed" -Context "Failed path task"
Assert-True ($failedTask.result) "Failed task did not include failure result"
$report.Add("worker failed path ok")

Write-Output "LR2 verification passed:"
$report | ForEach-Object { Write-Output "- $_" }
