export const normalizeList = (value) => {
    return Array.isArray(value) ? value : [];
};

export const countTasksByStatus = (tasks, status) => {
    return tasks.filter((task) => task.status === status).length;
};

export const getTaskResult = (task) => {
    if (task.status === "Completed" || task.status === "Failed") {
        return task.result || "-";
    }

    return "Task is still being processed.";
};
