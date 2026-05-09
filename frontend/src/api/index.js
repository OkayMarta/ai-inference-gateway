import * as authApi from "./authApi";
import * as modelsApi from "./modelsApi";
import * as tasksApi from "./tasksApi";
import { getToken, setToken } from "./tokenStorage";

export const api = {
    getToken,
    setToken,
    ...authApi,
    ...modelsApi,
    ...tasksApi,
};
