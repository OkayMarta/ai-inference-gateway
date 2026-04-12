import { useEffect, useState } from "react";
import { api } from "../api/client";

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

export default function Dashboard() {
    const [users, setUsers] = useState([]);
    const [models, setModels] = useState([]);
    const [tasks, setTasks] = useState([]);

    const [selectedUser, setSelectedUser] = useState("");
    const [selectedModel, setSelectedModel] = useState("");
    const [prompt, setPrompt] = useState("");

    const [loading, setLoading] = useState(false);
    const [loadingTasks, setLoadingTasks] = useState(false);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");

    useEffect(() => {
        api.getUsers().then(setUsers).catch((err) => setError(err.message));
        api.getModels().then(setModels).catch((err) => setError(err.message));
    }, []);

    useEffect(() => {
        if (!selectedUser) {
            setTasks([]);
            setSuccess("");
            return;
        }

        let active = true;

        const fetchTasks = async () => {
            setLoadingTasks(true);
            try {
                const [nextTasks, nextUsers] = await Promise.all([
                    api.getTasks(selectedUser),
                    api.getUsers(),
                ]);

                if (!active) {
                    return;
                }

                setTasks(nextTasks);
                setUsers(nextUsers);
            } catch (err) {
                if (active) {
                    setError(err.message);
                }
            } finally {
                if (active) {
                    setLoadingTasks(false);
                }
            }
        };

        fetchTasks();
        const intervalId = setInterval(fetchTasks, 2000);

        return () => {
            active = false;
            clearInterval(intervalId);
        };
    }, [selectedUser]);

    const handleSubmit = async (event) => {
        event.preventDefault();
        if (!selectedUser || !selectedModel || !prompt.trim()) {
            return;
        }

        setLoading(true);
        setError("");
        setSuccess("");

        try {
            await api.submitTask(selectedUser, selectedModel, prompt.trim());
            setPrompt("");
            setSuccess("Task submitted.");

            const nextTasks = await api.getTasks(selectedUser);
            setTasks(nextTasks);
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    const currentUser = users.find((user) => user.id === selectedUser);
    const currentModel = models.find((model) => model.id === selectedModel);
    const sortedTasks = [...tasks].sort(
        (left, right) => new Date(right.createdAt) - new Date(left.createdAt),
    );

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
                            onChange={(event) => {
                                setSelectedUser(event.target.value);
                                setError("");
                                setSuccess("");
                            }}
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
                            onChange={(event) => {
                                setSelectedModel(event.target.value);
                                setError("");
                                setSuccess("");
                            }}
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
                            disabled={!selectedUser || !selectedModel || loading}
                        />
                    </section>

                    {(error || success) && (
                        <div
                            className={`notice ${error ? "notice-error" : "notice-success"}`}
                        >
                            {error || success}
                        </div>
                    )}

                    <button
                        type="submit"
                        className="submit-button"
                        disabled={
                            loading ||
                            !selectedUser ||
                            !selectedModel ||
                            !prompt.trim()
                        }
                    >
                        {loading ? "Submitting..." : "Submit"}
                    </button>
                </form>
            </aside>

            <section className="panel task-panel">
                <div className="task-panel-header">
                    <div className="task-panel-meta">
                        <span>Tasks</span>
                        <span>{sortedTasks.length}</span>
                    </div>
                    {loadingTasks && (
                        <span className="task-panel-loading">Refreshing...</span>
                    )}
                </div>

                {!selectedUser ? (
                    <div className="empty-state">
                        Select a user to view task history.
                    </div>
                ) : sortedTasks.length === 0 ? (
                    <div className="empty-state">No tasks found for this user.</div>
                ) : (
                    <div className="task-list">
                        {sortedTasks.map((task) => {
                            const model = models.find(
                                (item) => item.id === task.modelId,
                            );

                            return (
                                <article key={task.id} className="task-card">
                                    <div className="task-card-header">
                                        <div className="task-title-row">
                                            <span className={getStatusClass(task.status)}>
                                                {task.status}
                                            </span>
                                            <span className="task-model">
                                                {model?.name || task.modelId}
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
                                            <dd>{model?.name || task.modelId}</dd>
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

