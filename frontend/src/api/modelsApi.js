import { request } from "./http";

export const getModels = () => request("/api/models");
