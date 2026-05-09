const TOKEN_STORAGE_KEY = "authToken";

export const getToken = () => localStorage.getItem(TOKEN_STORAGE_KEY);

export const setToken = (token) => {
    if (token) {
        localStorage.setItem(TOKEN_STORAGE_KEY, token);
    } else {
        localStorage.removeItem(TOKEN_STORAGE_KEY);
    }
};

export const clearToken = () => {
    localStorage.removeItem(TOKEN_STORAGE_KEY);
};
