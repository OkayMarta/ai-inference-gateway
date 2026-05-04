const API_BASE = "/api";
const DEFAULT_ERROR_MESSAGE = "Server request failed";
const TOKEN_STORAGE_KEY = "authToken";

const parseResponse = async (response) => {
    const contentType = response.headers.get("content-type") || "";
    const isJsonResponse = contentType.includes("application/json");

    if (isJsonResponse) {
        return response.json();
    }

    const text = await response.text();
    return text || null;
};

const getErrorMessage = (responseBody) => {
    if (responseBody && typeof responseBody === "object" && responseBody.message) {
        return responseBody.message;
    }

    if (typeof responseBody === "string" && responseBody.trim()) {
        return responseBody;
    }

    return DEFAULT_ERROR_MESSAGE;
};

const getToken = () => localStorage.getItem(TOKEN_STORAGE_KEY);

const setToken = (token) => {
    if (token) {
        localStorage.setItem(TOKEN_STORAGE_KEY, token);
    } else {
        localStorage.removeItem(TOKEN_STORAGE_KEY);
    }
};

const request = async (url, options = {}) => {
    const token = getToken();
    const headers = {
        "Content-Type": "application/json",
        ...options.headers,
    };

    if (token) {
        headers.Authorization = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE}${url}`, {
        headers,
        ...options,
    });

    const data = await parseResponse(response);

    if (!response.ok) {
        throw new Error(getErrorMessage(data));
    }

    return data;
};

export const api = {
    getToken,
    setToken,
    logout: () => setToken(""),
    login: async (email, password) => {
        const data = await request("/auth/login", {
            method: "POST",
            body: JSON.stringify({ email, password }),
        });
        setToken(data.token);
        return data;
    },
    register: async (username, email, password) => {
        const data = await request("/auth/register", {
            method: "POST",
            body: JSON.stringify({ username, email, password }),
        });
        setToken(data.token);
        return data;
    },
    getMe: () => request("/auth/me"),
    getUsers: () => request("/users"),
    getModels: () => request("/models"),
    getTasks: (params = {}) => {
        const query = new URLSearchParams();

        if (params.userId) {
            query.set("userId", params.userId);
        }
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
        return request(`/tasks${suffix}`);
    },
    submitTask: (modelId, payload) => {
        return request("/tasks", {
            method: "POST",
            body: JSON.stringify({ modelId, payload }),
        });
    },
    deleteTask: (taskId) =>
        request(`/tasks/${encodeURIComponent(taskId)}`, {
            method: "DELETE",
        }),
};
