const TaskSummary = ({
    queuedCount,
    processingCount,
    completedCount,
    failedCount,
    cancelledCount,
}) => (
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

export default TaskSummary;
