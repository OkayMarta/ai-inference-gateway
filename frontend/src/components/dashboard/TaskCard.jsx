import StatusBadge from "../StatusBadge";
import { formatTimestamp } from "../../utils/dateUtils";
import { getTaskResult } from "../../utils/taskUtils";

const TaskCard = ({
    task,
    models,
    isExpanded,
    onToggleExpanded,
    onCancelTask,
    cancelLoadingTaskId,
}) => {
    const taskModel =
        models.find((model) => model.id === task.modelId) || null;
    const taskResult = getTaskResult(task);
    const isExpandableResult = taskResult.length > 220;

    return (
        <article className="task-card">
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
                    onClick={onToggleExpanded}
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
};

export default TaskCard;
