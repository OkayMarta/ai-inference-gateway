const API_BASE = "/api";
const DEFAULT_ERROR_MESSAGE = "Server request failed";

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

const normalizeSubmitTaskPayload = (userIdOrParams, modelId, payload) => {
    if (userIdOrParams && typeof userIdOrParams === "object") {
        return userIdOrParams;
    }

    return { userId: userIdOrParams, modelId, payload };
};

const request = async (url, options = {}) => {
    const response = await fetch(`${API_BASE}${url}`, {
        headers: {
            "Content-Type": "application/json",
            ...options.headers,
        },
        ...options,
    });

    const data = await parseResponse(response);

    if (!response.ok) {
        throw new Error(getErrorMessage(data));
    }

    return data;
};

export const api = {
    getUsers: () => request("/users"),
    getModels: () => request("/models"),
    getTasks: (userId) => request(`/tasks?userId=${encodeURIComponent(userId)}`),
    submitTask: (userIdOrParams, modelId, payload) => {
        const task = normalizeSubmitTaskPayload(
            userIdOrParams,
            modelId,
            payload,
        );

        return request("/tasks", {
            method: "POST",
            body: JSON.stringify(task),
        });
    },
};
