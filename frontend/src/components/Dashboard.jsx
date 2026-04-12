import { useEffect, useMemo, useState } from "react";
import { api } from "../api/client";

function normalizeList(value) {
    return Array.isArray(value) ? value : [];
}

function formatTimestamp(value) {
    if (!value) {
        return "-";
    }

    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
        return value;
    }

    return new Intl.DateTimeFormat("uk-UA", {
        year: "numeric",
        month: "2-digit",
        day: "2-digit",
        hour: "2-digit",
        minute: "2-digit",
    }).format(date);
}

function getStatusClass(status) {
    switch (status) {
        case "Completed":
            return "status-badge status-completed";
        case "Processing":
            return "status-badge status-processing";
        case "Failed":
            return "status-badge status-failed";
        default:
            return "status-badge status-queued";
    }
}

function getTaskResult(task) {
    if (task.status === "Completed" || task.status === "Failed") {
        return task.result || "-";
    }

    return "Task is still being processed.";
}

function countTasksByStatus(tasks, status) {
    return tasks.filter((task) => task.status === status).length;
}

export default function Dashboard() {
    const [users, setUsers] = useState([]);
    const [models, setModels] = useState([]);
    const [tasks, setTasks] = useState([]);

    const [selectedUser, setSelectedUser] = useState("");
    const [selectedModel, setSelectedModel] = useState("");
    const [prompt, setPrompt] = useState("");

    const [bootLoading, setBootLoading] = useState(true);
    const [taskLoading, setTaskLoading] = useState(false);
    const [submitLoading, setSubmitLoading] = useState(false);
    const [screenError, setScreenError] = useState("");
    const [submitError, setSubmitError] = useState("");
    const [submitSuccess, setSubmitSuccess] = useState("");

    const currentUser = useMemo(
        () => users.find((user) => user.id === selectedUser) || null,
        [users, selectedUser],
    );

    const currentModel = useMemo(
        () => models.find((model) => model.id === selectedModel) || null,
        [models, selectedModel],
    );

    const sortedTasks = useMemo(
        () =>
            normalizeList(tasks).sort(
                (left, right) =>
                    new Date(right.createdAt) - new Date(left.createdAt),
            ),
        [tasks],
    );

    const queuedCount = useMemo(
        () => countTasksByStatus(sortedTasks, "Queued"),
        [sortedTasks],
    );
    const processingCount = useMemo(
        () => countTasksByStatus(sortedTasks, "Processing"),
        [sortedTasks],
    );
    const completedCount = useMemo(
        () => countTasksByStatus(sortedTasks, "Completed"),
        [sortedTasks],
    );

    useEffect(() => {
        let active = true;

        const loadBootData = async () => {
            setBootLoading(true);
            setScreenError("");

            try {
                const [nextUsers, nextModels] = await Promise.all([
                    api.getUsers(),
                    api.getModels(),
                ]);

                if (!active) {
                    return;
                }

                setUsers(normalizeList(nextUsers));
                setModels(normalizeList(nextModels));
            } catch (error) {
                if (active) {
                    setScreenError(error.message);
                }
            } finally {
                if (active) {
                    setBootLoading(false);
                }
            }
        };

        loadBootData();

        return () => {
            active = false;
        };
    }, []);

    useEffect(() => {
        if (!selectedUser) {
            setTasks([]);
            setTaskLoading(false);
            return;
        }

        let active = true;

        const refreshUserData = async () => {
            setTaskLoading(true);
            setScreenError("");

            try {
                const [nextTasks, nextUsers] = await Promise.all([
                    api.getTasks(selectedUser),
                    api.getUsers(),
                ]);

                if (!active) {
                    return;
                }

                setTasks(normalizeList(nextTasks));
                setUsers(normalizeList(nextUsers));
            } catch (error) {
                if (active) {
                    setScreenError(error.message);
                }
            } finally {
                if (active) {
                    setTaskLoading(false);
                }
            }
        };

        refreshUserData();
        const intervalId = setInterval(refreshUserData, 2000);

        return () => {
            active = false;
            clearInterval(intervalId);
        };
    }, [selectedUser]);

    const handleUserChange = (event) => {
        setSelectedUser(event.target.value);
        setSubmitError("");
        setSubmitSuccess("");
        setScreenError("");
    };

    const handleModelChange = (event) => {
        setSelectedModel(event.target.value);
        setSubmitError("");
        setSubmitSuccess("");
    };

    const handleSubmit = async (event) => {
        event.preventDefault();
        if (!selectedUser || !selectedModel || !prompt.trim()) {
            return;
        }

        setSubmitLoading(true);
        setSubmitError("");
        setSubmitSuccess("");

        try {
            await api.submitTask(selectedUser, selectedModel, prompt.trim());

            const [nextTasks, nextUsers] = await Promise.all([
                api.getTasks(selectedUser),
                api.getUsers(),
            ]);

            setTasks(normalizeList(nextTasks));
            setUsers(normalizeList(nextUsers));
            setPrompt("");
            setSubmitSuccess("Task submitted.");
        } catch (error) {
            setSubmitError(error.message);
        } finally {
            setSubmitLoading(false);
        }
    };

    if (bootLoading) {
        return (
            <div className="dashboard-layout">
                <aside className="panel control-panel">
                    <div className="empty-state">Loading dashboard...</div>
                </aside>
                <section className="panel task-panel">
                    <div className="empty-state">Loading tasks...</div>
                </section>
            </div>
        );
    }

    return (
        <div className="dashboard-layout">
            <aside className="panel control-panel">
                <form className="control-form" onSubmit={handleSubmit}>
                    <section className="control-section">
                        <label className="field-label" htmlFor="user-select">
                            User
                        </label>
                        <select
                            id="user-select"
                            value={selectedUser}
                            onChange={handleUserChange}
                            className="field-input"
                        >
                            <option value="">Select user</option>
                            {users.map((user) => (
                                <option key={user.id} value={user.id}>
                                    {user.username}
                                </option>
                            ))}
                        </select>
                    </section>

                    <section className="control-section">
                        <label className="field-label" htmlFor="model-select">
                            Model
                        </label>
                        <select
                            id="model-select"
                            value={selectedModel}
                            onChange={handleModelChange}
                            className="field-input"
                        >
                            <option value="">Select model</option>
                            {models.map((model) => (
                                <option key={model.id} value={model.id}>
                                    {model.name}
                                </option>
                            ))}
                        </select>
                    </section>

                    <section className="metrics-grid">
                        <div className="metric-card">
                            <span className="metric-label">Balance</span>
                            <span className="metric-value">
                                {currentUser
                                    ? currentUser.tokenBalance.toFixed(1)
                                    : "-"}
                            </span>
                        </div>
                        <div className="metric-card">
                            <span className="metric-label">Model cost</span>
                            <span className="metric-value">
                                {currentModel
                                    ? currentModel.tokenCost.toFixed(1)
                                    : "-"}
                            </span>
                        </div>
                    </section>

                    <section className="control-section">
                        <label className="field-label" htmlFor="prompt-input">
                            Prompt
                        </label>
                        <textarea
                            id="prompt-input"
                            value={prompt}
                            onChange={(event) => setPrompt(event.target.value)}
                            className="field-input field-textarea"
                            placeholder="Enter prompt"
                            disabled={
                                !selectedUser || !selectedModel || submitLoading
                            }
                        />
                    </section>

                    {screenError && (
                        <div className="notice notice-error">{screenError}</div>
                    )}
                    {submitError && (
                        <div className="notice notice-error">{submitError}</div>
                    )}
                    {submitSuccess && (
                        <div className="notice notice-success">
                            {submitSuccess}
                        </div>
                    )}

                    <button
                        type="submit"
                        className="submit-button"
                        disabled={
                            submitLoading ||
                            !selectedUser ||
                            !selectedModel ||
                            !prompt.trim()
                        }
                    >
                        {submitLoading ? "Submitting..." : "Submit"}
                    </button>
                </form>
            </aside>

            <section className="panel task-panel">
                <div className="task-panel-header">
                    <div className="task-panel-meta">
                        <span>Tasks</span>
                        <span>{sortedTasks.length}</span>
                    </div>
                    <div className="task-summary">
                        <span className="task-summary-item">
                            Queued {queuedCount}
                        </span>
                        <span className="task-summary-item">
                            Processing {processingCount}
                        </span>
                        <span className="task-summary-item">
                            Completed {completedCount}
                        </span>
                    </div>
                </div>

                {screenError && !selectedUser ? (
                    <div className="empty-state">{screenError}</div>
                ) : !selectedUser ? (
                    <div className="empty-state">
                        Select a user to view tasks.
                    </div>
                ) : taskLoading && sortedTasks.length === 0 ? (
                    <div className="empty-state">Loading tasks...</div>
                ) : sortedTasks.length === 0 ? (
                    <div className="empty-state">
                        No tasks available for the selected user.
                    </div>
                ) : (
                    <div className="task-list">
                        {sortedTasks.map((task) => {
                            const taskModel =
                                models.find((model) => model.id === task.modelId) ||
                                null;

                            return (
                                <article key={task.id} className="task-card">
                                    <div className="task-card-header">
                                        <div className="task-title-row">
                                            <span className={getStatusClass(task.status)}>
                                                {task.status}
                                            </span>
                                            <span className="task-model">
                                                {taskModel?.name || task.modelId}
                                            </span>
                                        </div>
                                        <span className="task-created-at">
                                            {formatTimestamp(task.createdAt)}
                                        </span>
                                    </div>

                                    <dl className="task-details-grid">
                                        <div className="task-detail">
                                            <dt>Task ID</dt>
                                            <dd>{task.id}</dd>
                                        </div>
                                        <div className="task-detail">
                                            <dt>Model</dt>
                                            <dd>
                                                {taskModel?.name || task.modelId}
                                            </dd>
                                        </div>
                                        <div className="task-detail task-detail-wide">
                                            <dt>Prompt</dt>
                                            <dd>{task.payload}</dd>
                                        </div>
                                        <div className="task-detail task-detail-wide">
                                            <dt>Result</dt>
                                            <dd>{getTaskResult(task)}</dd>
                                        </div>
                                    </dl>
                                </article>
                            );
                        })}
                    </div>
                )}
            </section>
        </div>
    );
}
