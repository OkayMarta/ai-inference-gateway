import { useEffect, useMemo, useState } from "react";
import { api } from "../api/client";
import AuthForm from "./AuthForm";
import EmptyState from "./EmptyState";
import LandingPage from "./landing/LandingPage";
import DashboardDrawer from "./layout/DashboardDrawer";
import MobileDashboardBar from "./layout/MobileDashboardBar";
import SessionBar from "./layout/SessionBar";
import SectionCard from "./SectionCard";
import TaskComposer from "./TaskComposer";
import TaskList from "./TaskList";
import { countTasksByStatus, normalizeList } from "../utils/taskUtils";
import { getInitials, sameUserSnapshot } from "../utils/userUtils";
import "../styles/components/Dashboard.css";

const Dashboard = () => {
    const [passwordResetToken, setPasswordResetToken] = useState(() => {
        return new URLSearchParams(window.location.search).get("token") || "";
    });
    const [landingStarted, setLandingStarted] = useState(false);
    const [authUser, setAuthUser] = useState(null);
    const [currentUser, setCurrentUser] = useState(null);
    const [models, setModels] = useState([]);
    const [tasks, setTasks] = useState([]);

    const [selectedModelId, setSelectedModelId] = useState("");
    const [prompt, setPrompt] = useState("");
    const [statusFilter, setStatusFilter] = useState("");
    const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

    const [authLoading, setAuthLoading] = useState(Boolean(api.getToken()));
    const [authError, setAuthError] = useState("");
    const [authSuccess, setAuthSuccess] = useState("");
    const [bootLoading, setBootLoading] = useState(false);
    const [taskLoading, setTaskLoading] = useState(false);
    const [submitLoading, setSubmitLoading] = useState(false);
    const [cancelLoadingTaskId, setCancelLoadingTaskId] = useState("");
    const [screenError, setScreenError] = useState("");
    const [submitError, setSubmitError] = useState("");
    const [submitSuccess, setSubmitSuccess] = useState("");
    const [balanceAlert, setBalanceAlert] = useState("");
    const [metricFlashToken, setMetricFlashToken] = useState(0);

    const currentModel = useMemo(
        () => models.find((model) => model.id === selectedModelId) || null,
        [models, selectedModelId],
    );
    const hasAvailableModels = models.length > 0;

    const sortedTasks = useMemo(
        () =>
            [...normalizeList(tasks)].sort(
                (left, right) =>
                    new Date(right.createdAt) - new Date(left.createdAt),
            ),
        [tasks],
    );

    const queuedCount = useMemo(
        () => countTasksByStatus(sortedTasks, "Queued"),
        [sortedTasks],
    );
    const processingCount = useMemo(
        () => countTasksByStatus(sortedTasks, "Processing"),
        [sortedTasks],
    );
    const completedCount = useMemo(
        () => countTasksByStatus(sortedTasks, "Completed"),
        [sortedTasks],
    );
    const failedCount = useMemo(
        () => countTasksByStatus(sortedTasks, "Failed"),
        [sortedTasks],
    );
    const cancelledCount = useMemo(
        () => countTasksByStatus(sortedTasks, "Cancelled"),
        [sortedTasks],
    );

    const hydrateCurrentUser = async (me) => {
        const userID = me.id || me.userId || me.ID;
        const users = normalizeList(await api.getUsers());
        return users.find((user) => user.id === userID) || {
            id: userID,
            email: me.email || me.Email,
            role: me.role,
        };
    };

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
    }, [passwordResetToken]);

    useEffect(() => {
        if (!currentUser) {
            return undefined;
        }

        let active = true;

        const loadBootData = async () => {
            setBootLoading(true);
            setScreenError("");

            try {
                const nextModels = await api.getModels();
                if (!active) {
                    return;
                }
                setModels(normalizeList(nextModels));
            } catch (error) {
                if (active) {
                    setScreenError(error.message);
                }
            } finally {
                if (active) {
                    setBootLoading(false);
                }
            }
        };

        loadBootData();

        return () => {
            active = false;
        };
    }, [currentUser?.id]);

    useEffect(() => {
        if (!currentUser) {
            setTasks([]);
            setTaskLoading(false);
            return undefined;
        }

        let active = true;

        const refreshUserData = async (showLoadingState = false) => {
            if (showLoadingState) {
                setTaskLoading(true);
            }
            setScreenError("");

            try {
                const [nextTasks, nextUsers] = await Promise.all([
                    api.getTasks({
                        userId: currentUser.id,
                        status: statusFilter,
                        limit: 20,
                        offset: 0,
                        sort: "created_at_desc",
                    }),
                    api.getUsers(),
                ]);

                if (!active) {
                    return;
                }

                setTasks(normalizeList(nextTasks));
                const refreshedUser = normalizeList(nextUsers).find(
                    (user) => user.id === currentUser.id,
                );
                if (refreshedUser) {
                    setCurrentUser((previousUser) =>
                        sameUserSnapshot(previousUser, refreshedUser)
                            ? previousUser
                            : refreshedUser,
                    );
                }
            } catch (error) {
                if (active) {
                    setScreenError(error.message);
                }
            } finally {
                if (active && showLoadingState) {
                    setTaskLoading(false);
                }
            }
        };

        refreshUserData(true);
        const intervalId = setInterval(() => {
            refreshUserData(false);
        }, 2000);

        return () => {
            active = false;
            clearInterval(intervalId);
        };
    }, [currentUser?.id, statusFilter]);

    useEffect(() => {
        if (!submitError) {
            return undefined;
        }

        const timeoutId = window.setTimeout(() => {
            setSubmitError("");
        }, 4200);

        return () => window.clearTimeout(timeoutId);
    }, [submitError]);

    useEffect(() => {
        if (!submitSuccess) {
            return undefined;
        }

        const timeoutId = window.setTimeout(() => {
            setSubmitSuccess("");
        }, 2400);

        return () => window.clearTimeout(timeoutId);
    }, [submitSuccess]);

    useEffect(() => {
        if (!balanceAlert) {
            return undefined;
        }

        const timeoutId = window.setTimeout(() => {
            setBalanceAlert("");
        }, 4200);

        return () => window.clearTimeout(timeoutId);
    }, [balanceAlert]);

    useEffect(() => {
        if (!mobileMenuOpen) {
            return undefined;
        }

        const handleKeyDown = (event) => {
            if (event.key === "Escape") {
                setMobileMenuOpen(false);
            }
        };

        window.addEventListener("keydown", handleKeyDown);

        return () => window.removeEventListener("keydown", handleKeyDown);
    }, [mobileMenuOpen]);

    const completeAuth = async (authAction) => {
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
    };

    const handleLogin = ({ email, password }) => {
        completeAuth(() => api.login(email, password));
    };

    const handleRegister = ({ username, email, password }) => {
        completeAuth(() => api.register(username, email, password));
    };

    const handleForgotPassword = async ({ email }) => {
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
    };

    const handleResetPassword = async ({ token, newPassword }) => {
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
    };

    const handleLogout = () => {
        api.logout();
        setAuthUser(null);
        setCurrentUser(null);
        setModels([]);
        setTasks([]);
        setSelectedModelId("");
        setPrompt("");
        setStatusFilter("");
        setScreenError("");
        setAuthError("");
        setAuthSuccess("");
        setLandingStarted(false);
    };

    const handleModelChange = (event) => {
        setSelectedModelId(event.target.value);
        setSubmitError("");
        setSubmitSuccess("");
        setBalanceAlert("");
    };

    const handlePromptChange = (event) => {
        setSubmitError("");
        setSubmitSuccess("");
        setPrompt(event.target.value);
    };

    const handleStatusFilterChange = (event) => {
        setStatusFilter(event.target.value);
        setScreenError("");
    };

    const handleBackToLanding = () => {
        setLandingStarted(false);
        window.history.replaceState({}, "", window.location.pathname);
    };

    const handleSubmit = async (event) => {
        event.preventDefault();
        if (!selectedModelId || !prompt.trim()) {
            return;
        }

        setSubmitLoading(true);
        setSubmitError("");
        setSubmitSuccess("");
        setBalanceAlert("");

        try {
            await api.submitTask(selectedModelId, prompt.trim());

            const [nextTasks, nextUsers] = await Promise.all([
                api.getTasks({
                    userId: currentUser.id,
                    status: statusFilter,
                    limit: 20,
                    offset: 0,
                    sort: "created_at_desc",
                }),
                api.getUsers(),
            ]);

            setTasks(normalizeList(nextTasks));
            const refreshedUser = normalizeList(nextUsers).find(
                (user) => user.id === currentUser.id,
            );
            if (refreshedUser) {
                setCurrentUser((previousUser) =>
                    sameUserSnapshot(previousUser, refreshedUser)
                        ? previousUser
                        : refreshedUser,
                );
            }
            setPrompt("");
            setSubmitSuccess("Task submitted.");
        } catch (error) {
            if (error.message === "insufficient token balance") {
                setBalanceAlert("Insufficient token balance.");
                setMetricFlashToken((value) => value + 1);
            } else {
                setSubmitError(error.message);
            }
        } finally {
            setSubmitLoading(false);
        }
    };

    const handleCancelTask = async (taskId) => {
        try {
            setCancelLoadingTaskId(taskId);
            await api.deleteTask(taskId);

            const [nextTasks, nextUsers] = await Promise.all([
                api.getTasks({
                    userId: currentUser.id,
                    status: statusFilter,
                    limit: 20,
                    offset: 0,
                    sort: "created_at_desc",
                }),
                api.getUsers(),
            ]);

            setTasks(normalizeList(nextTasks));
            const refreshedUser = normalizeList(nextUsers).find(
                (user) => user.id === currentUser.id,
            );
            if (refreshedUser) {
                setCurrentUser((previousUser) =>
                    sameUserSnapshot(previousUser, refreshedUser)
                        ? previousUser
                        : refreshedUser,
                );
            }
        } catch (error) {
            console.error(error);
            window.alert(`Error: ${error.message}`);
        } finally {
            setCancelLoadingTaskId("");
        }
    };

    if (authLoading) {
        return (
            <div className="auth-shell">
                <SectionCard className="auth-card">
                    <EmptyState
                        title="Checking session"
                        description="Authentication state is being restored."
                    />
                </SectionCard>
            </div>
        );
    }

    if (!landingStarted && !api.getToken()) {
        return (
            <LandingPage
                onStart={() => setLandingStarted(true)}
                onLogin={() => setLandingStarted(true)}
            />
        );
    }

    if (!currentUser) {
        return (
            <AuthForm
                onLogin={handleLogin}
                onRegister={handleRegister}
                onForgotPassword={handleForgotPassword}
                onResetPassword={handleResetPassword}
                onBackToLanding={handleBackToLanding}
                loading={authLoading}
                error={authError}
                success={authSuccess}
                resetToken={passwordResetToken}
            />
        );
    }

    if (bootLoading) {
        return (
            <div className="dashboard-page">
                <div className="dashboard-glow dashboard-glow-primary" aria-hidden="true" />
                <div className="dashboard-glow dashboard-glow-secondary" aria-hidden="true" />

                <div className="dashboard-stack">
                    <div className="dashboard-layout dashboard-layout-loading">
                        <SectionCard as="aside" className="control-panel">
                            <EmptyState
                                title="Loading dashboard"
                                description="Initial data is being loaded."
                            />
                        </SectionCard>
                        <SectionCard className="task-panel">
                            <EmptyState
                                title="Loading tasks"
                                description="Task history is being prepared."
                            />
                        </SectionCard>
                    </div>
                </div>
            </div>
        );
    }

    const displayName =
        currentUser.username || authUser?.email || currentUser.email || "Account";
    const displayEmail = currentUser.email || authUser?.email || "";
    const initials = getInitials(displayName || displayEmail);
    const sidebarId = "dashboard-sidebar";

    return (
        <div
            className={`dashboard-page${mobileMenuOpen ? " dashboard-menu-open" : ""}`}
            id="top"
        >
            <div className="dashboard-glow dashboard-glow-primary" aria-hidden="true" />
            <div className="dashboard-glow dashboard-glow-secondary" aria-hidden="true" />

            <div className="dashboard-stack">
                <MobileDashboardBar
                    mobileMenuOpen={mobileMenuOpen}
                    onOpenMenu={() => setMobileMenuOpen(true)}
                    sidebarId={sidebarId}
                />

                <SessionBar
                    displayEmail={displayEmail}
                    displayName={displayName}
                    initials={initials}
                    onLogout={handleLogout}
                />

                <div className="dashboard-layout">
                    <DashboardDrawer
                        mobileMenuOpen={mobileMenuOpen}
                        onCloseMenu={() => setMobileMenuOpen(false)}
                        sidebarId={sidebarId}
                    >
                        <TaskComposer
                            models={models}
                            hasAvailableModels={hasAvailableModels}
                            selectedModelId={selectedModelId}
                            prompt={prompt}
                            currentUser={currentUser}
                            currentModel={currentModel}
                            screenError={screenError}
                            submitError={submitError}
                            submitSuccess={submitSuccess}
                            balanceAlert={balanceAlert}
                            metricFlashToken={metricFlashToken}
                            submitLoading={submitLoading}
                            onDismissSubmitError={() => setSubmitError("")}
                            onDismissSubmitSuccess={() => setSubmitSuccess("")}
                            onDismissBalanceAlert={() => setBalanceAlert("")}
                            onModelChange={handleModelChange}
                            onPromptChange={handlePromptChange}
                            onSubmit={handleSubmit}
                        />

                        <SessionBar
                            displayEmail={displayEmail}
                            displayName={displayName}
                            initials={initials}
                            onLogout={handleLogout}
                            variant="drawer"
                        />
                    </DashboardDrawer>

                    <TaskList
                        models={models}
                        selectedUserId={currentUser.id}
                        screenError={screenError}
                        taskLoading={taskLoading}
                        sortedTasks={sortedTasks}
                        queuedCount={queuedCount}
                        processingCount={processingCount}
                        completedCount={completedCount}
                        failedCount={failedCount}
                        cancelledCount={cancelledCount}
                        statusFilter={statusFilter}
                        onStatusFilterChange={handleStatusFilterChange}
                        onCancelTask={handleCancelTask}
                        cancelLoadingTaskId={cancelLoadingTaskId}
                    />
                </div>
            </div>
        </div>
    );
};

export default Dashboard;
