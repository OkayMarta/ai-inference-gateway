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
}) => {
    const taskSummary = (
        <div className="task-summary">
            <span className="task-summary-item">Queued {queuedCount}</span>
            <span className="task-summary-item">Processing {processingCount}</span>
            <span className="task-summary-item">Completed {completedCount}</span>
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
    } else if (taskLoading && sortedTasks.length === 0) {
        content = (
            <EmptyState title="Loading tasks" description="Task history is being loaded." />
        );
    } else if (screenError && sortedTasks.length === 0) {
        content = (
            <EmptyState title="Unable to load tasks" description={screenError} />
        );
    } else if (sortedTasks.length === 0) {
        content = (
            <EmptyState
                title="No tasks available"
                description="The selected user has no tasks yet."
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
            </div>
        );
    }

    return (
        <SectionCard
            className="task-panel"
            title="Tasks"
            rightSlot={
                <div className="task-panel-toolbar">
                    <div className="task-panel-count">{sortedTasks.length}</div>
                    {taskSummary}
                </div>
            }
        >
            {content}
        </SectionCard>
    );
};

export default TaskList;
