import { protectedRequest } from "./http";

export const getModels = () => protectedRequest("/api/models");
