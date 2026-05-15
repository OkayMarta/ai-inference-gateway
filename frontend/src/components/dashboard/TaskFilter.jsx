import { TASK_STATUSES } from "../../utils/taskUtils";

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
            {TASK_STATUSES.map((status) => (
                <option key={status} value={status}>
                    {status}
                </option>
            ))}
        </select>
    </div>
);

export default TaskFilter;
