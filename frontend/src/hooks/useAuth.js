import { useCallback, useEffect, useState } from "react";
import { api } from "../api/client";
import { normalizeList } from "../utils/taskUtils";

const getPasswordResetToken = () => {
    return new URLSearchParams(window.location.search).get("token") || "";
};

const useAuth = () => {
    const [passwordResetToken, setPasswordResetToken] = useState(getPasswordResetToken);
    const [landingStarted, setLandingStarted] = useState(false);
    const [authUser, setAuthUser] = useState(null);
    const [currentUser, setCurrentUser] = useState(null);
    const [authLoading, setAuthLoading] = useState(Boolean(api.getToken()));
    const [authError, setAuthError] = useState("");
    const [authSuccess, setAuthSuccess] = useState("");

    const hydrateCurrentUser = useCallback(async (me) => {
        const userID = me.id || me.userId || me.ID;
        const users = normalizeList(await api.getUsers());
        return users.find((user) => user.id === userID) || {
            id: userID,
            email: me.email || me.Email,
            role: me.role,
        };
    }, []);

    useEffect(() => {
        let active = true;

        const restoreSession = async () => {
            if (passwordResetToken) {
                api.logout();
                setAuthLoading(false);
                setLandingStarted(true);
                return;
            }

            if (!api.getToken()) {
                setAuthLoading(false);
                return;
            }

            try {
                const me = await api.getMe();
                const user = await hydrateCurrentUser(me);

                if (!active) {
                    return;
                }

                setAuthUser(me);
                setCurrentUser(user);
                setLandingStarted(true);
            } catch (error) {
                api.logout();
                if (active) {
                    setAuthError("Session expired. Please login again.");
                    setLandingStarted(true);
                }
            } finally {
                if (active) {
                    setAuthLoading(false);
                }
            }
        };

        restoreSession();

        return () => {
            active = false;
        };
    }, [hydrateCurrentUser, passwordResetToken]);

    const completeAuth = useCallback(
        async (authAction) => {
            setAuthLoading(true);
            setAuthError("");
            setAuthSuccess("");

            try {
                await authAction();
                const me = await api.getMe();
                const user = await hydrateCurrentUser(me);
                setAuthUser(me);
                setCurrentUser(user);
                setLandingStarted(true);
            } catch (error) {
                api.logout();
                setAuthError(error.message);
            } finally {
                setAuthLoading(false);
            }
        },
        [hydrateCurrentUser],
    );

    const login = useCallback(
        ({ email, password }) => {
            completeAuth(() => api.login(email, password));
        },
        [completeAuth],
    );

    const register = useCallback(
        ({ username, email, password }) => {
            completeAuth(() => api.register(username, email, password));
        },
        [completeAuth],
    );

    const forgotPassword = useCallback(async ({ email }) => {
        setAuthLoading(true);
        setAuthError("");
        setAuthSuccess("");

        try {
            const response = await api.requestPasswordReset(email);
            setAuthSuccess(response.message);
        } catch (error) {
            setAuthError(error.message);
        } finally {
            setAuthLoading(false);
        }
    }, []);

    const resetPassword = useCallback(async ({ token, newPassword }) => {
        setAuthLoading(true);
        setAuthError("");
        setAuthSuccess("");

        try {
            const response = await api.resetPassword(token, newPassword);
            setAuthSuccess(`${response.message} You can now log in.`);
            setPasswordResetToken("");
            window.history.replaceState({}, "", window.location.pathname);
        } catch (error) {
            setAuthError(error.message);
        } finally {
            setAuthLoading(false);
        }
    }, []);

    const logout = useCallback(() => {
        api.logout();
        setAuthUser(null);
        setCurrentUser(null);
        setAuthError("");
        setAuthSuccess("");
        setLandingStarted(false);
    }, []);

    const backToLanding = useCallback(() => {
        setLandingStarted(false);
        window.history.replaceState({}, "", window.location.pathname);
    }, []);

    return {
        authError,
        authLoading,
        authSuccess,
        authUser,
        backToLanding,
        currentUser,
        forgotPassword,
        hasAuthToken: Boolean(api.getToken()),
        landingStarted,
        login,
        logout,
        passwordResetToken,
        register,
        resetPassword,
        setCurrentUser,
        startLanding: () => setLandingStarted(true),
    };
};

export default useAuth;
