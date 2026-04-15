import EmptyState from "./EmptyState";
import SectionCard from "./SectionCard";
import StatusBadge from "./StatusBadge";

const formatTimestamp = (value) => {
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
};

const getTaskResult = (task) => {
    if (task.status === "Completed" || task.status === "Failed") {
        return task.result || "-";
    }

    return "Task is still being processed.";
};

const TaskList = ({
    models,
    selectedUserId,
    screenError,
    taskLoading,
    sortedTasks,
    queuedCount,
    processingCount,
    completedCount,
    failedCount,
    cancelledCount,
    statusFilter,
    onStatusFilterChange,
    onCancelTask,
    cancelLoadingTaskId,
}) => {
    const hasSelectedUser = Boolean(selectedUserId);

    const taskSummary = (
        <div className="task-summary">
            <span className="task-summary-item task-summary-queued">Queued {queuedCount}</span>
            <span className="task-summary-item task-summary-processing">
                Processing {processingCount}
            </span>
            <span className="task-summary-item task-summary-completed">
                Completed {completedCount}
            </span>
            <span className="task-summary-item task-summary-failed">Failed {failedCount}</span>
            <span className="task-summary-item task-summary-cancelled">
                Cancelled {cancelledCount}
            </span>
        </div>
    );
    const filterControls = (
        <div className="task-list-controls">
            <label className="field-label task-filter-label" htmlFor="task-status-filter">
                Status
            </label>
            <select
                id="task-status-filter"
                value={statusFilter}
                onChange={onStatusFilterChange}
                className="field-input task-filter-select"
            >
                <option value="">All</option>
                <option value="Queued">Queued</option>
                <option value="Processing">Processing</option>
                <option value="Completed">Completed</option>
                <option value="Failed">Failed</option>
                <option value="Cancelled">Cancelled</option>
            </select>
        </div>
    );

    let content = null;

    if (screenError && !selectedUserId) {
        content = (
            <EmptyState title={screenError} description="Try reloading the dashboard." />
        );
    } else if (!selectedUserId) {
        content = (
            <EmptyState
                title="No user selected"
                description="Select a user to view tasks."
            />
        );
    } else {
        content = (
            <div className="task-list-stack">
                {screenError && (
                    <div className="notice notice-error notice-quiet">
                        {screenError}
                    </div>
                )}
                {filterControls}
                {taskLoading && sortedTasks.length === 0 ? (
                    <EmptyState
                        title="Loading tasks"
                        description="Task history is being loaded."
                    />
                ) : sortedTasks.length === 0 ? (
                    <EmptyState
                        title="No tasks available"
                        description={
                            statusFilter
                                ? "No tasks match the selected status filter."
                                : "The selected user has no tasks yet."
                        }
                    />
                ) : (
                    <div className="task-list">
                        {sortedTasks.map((task) => {
                            const taskModel =
                                models.find((model) => model.id === task.modelId) || null;

                            return (
                                <article key={task.id} className="task-card">
                                    <div className="task-card-header">
                                        <div className="task-title-row">
                                            <StatusBadge status={task.status} />
                                            <span className="task-model">
                                                {taskModel?.name || task.modelId}
                                            </span>
                                        </div>
                                        <span className="task-created-at">
                                            {formatTimestamp(task.createdAt)}
                                        </span>
                                    </div>

                                    {task.status === "Queued" && (
                                        <div className="task-card-actions">
                                            <button
                                                type="button"
                                                className="task-action-button"
                                                onClick={() => onCancelTask(task.id)}
                                                disabled={cancelLoadingTaskId === task.id}
                                            >
                                                {cancelLoadingTaskId === task.id
                                                    ? "Cancelling..."
                                                    : "Cancel"}
                                            </button>
                                        </div>
                                    )}

                                    <dl className="task-details-grid">
                                        <div className="task-detail">
                                            <dt>Task ID</dt>
                                            <dd>{task.id}</dd>
                                        </div>
                                        <div className="task-detail">
                                            <dt>Model</dt>
                                            <dd>{taskModel?.name || task.modelId}</dd>
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
            </div>
        );
    }

    return (
        <SectionCard
            className="task-panel"
            title="Tasks"
            rightSlot={
                <div className="task-panel-toolbar">
                    {hasSelectedUser ? (
                        taskSummary
                    ) : (
                        <div className="task-panel-helper">
                            Select a user to see tasks
                        </div>
                    )}
                </div>
            }
        >
            {content}
        </SectionCard>
    );
};

export default TaskList;
