import { getToken } from "./tokenStorage";

export const API_BASE_URL = "http://localhost:8080";

const DEFAULT_ERROR_MESSAGE = "Server request failed";

export const parseResponse = async (response) => {
    const contentType = response.headers.get("content-type") || "";
    const isJsonResponse = contentType.includes("application/json");

    if (isJsonResponse) {
        return response.json();
    }

    const text = await response.text();
    return text || null;
};

export const getErrorMessage = (responseBody) => {
    if (responseBody && typeof responseBody === "object" && responseBody.message) {
        return responseBody.message;
    }

    if (typeof responseBody === "string" && responseBody.trim()) {
        return responseBody;
    }

    return DEFAULT_ERROR_MESSAGE;
};

export const request = async (path, options = {}) => {
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

export const protectedRequest = async (path, options = {}) => {
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
