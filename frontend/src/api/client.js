const API_BASE_URL = "http://localhost:8080";
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

const request = async (path, options = {}) => {
    const headers = {
        "Content-Type": "application/json",
        ...options.headers,
    };

    const response = await fetch(`${API_BASE_URL}${path}`, {
        headers,
        ...options,
    });

    const data = await parseResponse(response);

    if (!response.ok) {
        throw new Error(getErrorMessage(data));
    }

    return data;
};

const protectedRequest = async (path, options = {}) => {
    const token = getToken();
    const headers = {
        ...options.headers,
    };

    if (token) {
        headers.Authorization = `Bearer ${token}`;
    }

    return request(path, {
        ...options,
        headers,
    });
};

export const api = {
    getToken,
    setToken,
    logout: () => setToken(""),
    login: async (email, password) => {
        const data = await request("/api/auth/login", {
            method: "POST",
            body: JSON.stringify({ email, password }),
        });
        setToken(data.token);
        return data;
    },
    register: async (username, email, password) => {
        const data = await request("/api/auth/register", {
            method: "POST",
            body: JSON.stringify({ username, email, password }),
        });
        setToken(data.token);
        return data;
    },
    requestPasswordReset: (email) =>
        request("/api/auth/forgot-password", {
            method: "POST",
            body: JSON.stringify({ email }),
        }),
    resetPassword: (token, newPassword) =>
        request("/api/auth/reset-password", {
            method: "POST",
            body: JSON.stringify({ token, newPassword }),
        }),
    getMe: () => protectedRequest("/api/auth/me"),
    getUsers: async () => [await protectedRequest("/api/auth/me")],
    getModels: () => request("/api/models"),
    getTasks: (params = {}) => {
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
    },
    submitTask: (modelId, payload) => {
        return protectedRequest("/api/tasks", {
            method: "POST",
            body: JSON.stringify({ modelId, payload }),
        });
    },
    getTask: (taskId) => protectedRequest(`/api/tasks/${encodeURIComponent(taskId)}`),
    updateTask: (taskId, payload) =>
        protectedRequest(`/api/tasks/${encodeURIComponent(taskId)}`, {
            method: "PUT",
            body: JSON.stringify({ payload }),
        }),
    deleteTask: (taskId) =>
        protectedRequest(`/api/tasks/${encodeURIComponent(taskId)}`, {
            method: "DELETE",
        }),
};
