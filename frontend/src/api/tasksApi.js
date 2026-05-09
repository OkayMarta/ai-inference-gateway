import { protectedRequest } from "./http";

export const getTasks = (params = {}) => {
    const query = new URLSearchParams();

    if (params.status) {
        query.set("status", params.status);
    }
    if (params.limit) {
        query.set("limit", String(params.limit));
    }
    if (typeof params.offset === "number") {
        query.set("offset", String(params.offset));
    }
    if (params.sort) {
        query.set("sort", params.sort);
    }

    const suffix = query.toString() ? `?${query.toString()}` : "";
    return protectedRequest(`/api/tasks${suffix}`);
};

export const submitTask = (modelId, payload) => {
    return protectedRequest("/api/tasks", {
        method: "POST",
        body: JSON.stringify({ modelId, payload }),
    });
};

export const getTask = (taskId) =>
    protectedRequest(`/api/tasks/${encodeURIComponent(taskId)}`);

export const updateTask = (taskId, payload) =>
    protectedRequest(`/api/tasks/${encodeURIComponent(taskId)}`, {
        method: "PUT",
        body: JSON.stringify({ payload }),
    });

export const deleteTask = (taskId) =>
    protectedRequest(`/api/tasks/${encodeURIComponent(taskId)}`, {
        method: "DELETE",
    });
