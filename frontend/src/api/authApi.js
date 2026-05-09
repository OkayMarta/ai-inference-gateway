import { protectedRequest, request } from "./http";
import { clearToken, setToken } from "./tokenStorage";

export const login = async (email, password) => {
    const data = await request("/api/auth/login", {
        method: "POST",
        body: JSON.stringify({ email, password }),
    });
    setToken(data.token);
    return data;
};

export const register = async (username, email, password) => {
    const data = await request("/api/auth/register", {
        method: "POST",
        body: JSON.stringify({ username, email, password }),
    });
    setToken(data.token);
    return data;
};

export const requestPasswordReset = (email) =>
    request("/api/auth/forgot-password", {
        method: "POST",
        body: JSON.stringify({ email }),
    });

export const resetPassword = (token, newPassword) =>
    request("/api/auth/reset-password", {
        method: "POST",
        body: JSON.stringify({ token, newPassword }),
    });

export const getMe = () => protectedRequest("/api/auth/me");

export const getUsers = async () => [await protectedRequest("/api/auth/me")];

export const logout = () => clearToken();
