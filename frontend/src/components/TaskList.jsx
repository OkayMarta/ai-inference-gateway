import { useState } from "react";
import TaskCard from "./dashboard/TaskCard";
import TaskFilter from "./dashboard/TaskFilter";
import TaskSummary from "./dashboard/TaskSummary";
import EmptyState from "./EmptyState";
import SectionCard from "./SectionCard";

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
        <TaskSummary
            queuedCount={queuedCount}
            processingCount={processingCount}
            completedCount={completedCount}
            failedCount={failedCount}
            cancelledCount={cancelledCount}
        />
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
                <TaskFilter
                    statusFilter={statusFilter}
                    onStatusFilterChange={onStatusFilterChange}
                />
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
                        {sortedTasks.map((task) => (
                            <TaskCard
                                key={task.id}
                                task={task}
                                models={models}
                                isExpanded={expandedTaskIds.has(task.id)}
                                onToggleExpanded={() => toggleExpandedTask(task.id)}
                                onCancelTask={onCancelTask}
                                cancelLoadingTaskId={cancelLoadingTaskId}
                            />
                        ))}
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
