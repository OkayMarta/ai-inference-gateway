const TaskFilter = ({ statusFilter, onStatusFilterChange }) => (
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

export default TaskFilter;
