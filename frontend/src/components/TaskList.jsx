import { useState } from "react";
import EmptyState from "./EmptyState";
import SectionCard from "./SectionCard";
import StatusBadge from "./StatusBadge";
import { formatTimestamp } from "../utils/dateUtils";
import { getTaskResult } from "../utils/taskUtils";

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
    const [expandedTaskIds, setExpandedTaskIds] = useState(() => new Set());

    const toggleExpandedTask = (taskId) => {
        setExpandedTaskIds((previousIds) => {
            const nextIds = new Set(previousIds);
            if (nextIds.has(taskId)) {
                nextIds.delete(taskId);
            } else {
                nextIds.add(taskId);
            }
            return nextIds;
        });
    };

    const panelTitle = (
        <span className="panel-title-with-icon">
            <svg viewBox="0 0 24 24" aria-hidden="true">
                <path d="M9 6h11" />
                <path d="M9 12h11" />
                <path d="M9 18h11" />
                <path d="M4 6h.01" />
                <path d="M4 12h.01" />
                <path d="M4 18h.01" />
            </svg>
            <span>Tasks</span>
        </span>
    );

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
                            const taskResult = getTaskResult(task);
                            const isExpandableResult = taskResult.length > 220;
                            const isExpanded = expandedTaskIds.has(task.id);

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
                                            <svg viewBox="0 0 24 24" aria-hidden="true">
                                                <path d="M8 2v4" />
                                                <path d="M16 2v4" />
                                                <path d="M3 10h18" />
                                                <rect x="3" y="4" width="18" height="18" rx="2" />
                                            </svg>
                                            <span>{formatTimestamp(task.createdAt)}</span>
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
                                            <dd className={isExpanded ? "" : "task-result-preview"}>
                                                {taskResult}
                                            </dd>
                                        </div>
                                    </dl>

                                    {isExpandableResult && (
                                        <button
                                            type="button"
                                            className="task-result-button"
                                            onClick={() => toggleExpandedTask(task.id)}
                                        >
                                            <span>{isExpanded ? "Show less" : "View full result"}</span>
                                            <svg viewBox="0 0 24 24" aria-hidden="true">
                                                {isExpanded ? (
                                                    <>
                                                        <path d="M18 6 6 18" />
                                                        <path d="M6 6h12v12" />
                                                    </>
                                                ) : (
                                                    <>
                                                        <path d="M7 17 17 7" />
                                                        <path d="M7 7h10v10" />
                                                    </>
                                                )}
                                            </svg>
                                        </button>
                                    )}
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
            title={panelTitle}
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
