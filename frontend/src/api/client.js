const API_BASE = "/api";

async function request(url, options = {}) {
    const res = await fetch(`${API_BASE}${url}`, {
        headers: { "Content-Type": "application/json" },
        ...options,
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.message || "Помилка сервера");
    return data;
}

export const api = {
    getUsers: () => request("/users"),
    getUser: (id) => request(`/users/${id}`),
    getModels: () => request("/models"),
    submitTask: (userId, modelId, payload) =>
        request("/tasks", {
            method: "POST",
            body: JSON.stringify({ userId, modelId, payload }),
        }),
    getTasks: (userId) => request(`/tasks?userId=${userId}`),
};
