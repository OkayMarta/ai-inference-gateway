import { useEffect, useMemo, useState } from "react";
import { api } from "../api/client";
import AuthForm from "./AuthForm";
import EmptyState from "./EmptyState";
import SectionCard from "./SectionCard";
import TaskComposer from "./TaskComposer";
import TaskList from "./TaskList";
import "../styles/components/Dashboard.css";

const normalizeList = (value) => {
    return Array.isArray(value) ? value : [];
};

const countTasksByStatus = (tasks, status) => {
    return tasks.filter((task) => task.status === status).length;
};

const sameUserSnapshot = (left, right) => {
    if (!left || !right) {
        return false;
    }

    return (
        left.id === right.id &&
        left.username === right.username &&
        left.email === right.email &&
        left.role === right.role &&
        left.tokenBalance === right.tokenBalance
    );
};

const Dashboard = () => {
    const [landingStarted, setLandingStarted] = useState(false);
    const [authUser, setAuthUser] = useState(null);
    const [currentUser, setCurrentUser] = useState(null);
    const [models, setModels] = useState([]);
    const [tasks, setTasks] = useState([]);

    const [selectedModelId, setSelectedModelId] = useState("");
    const [prompt, setPrompt] = useState("");
    const [statusFilter, setStatusFilter] = useState("");

    const [authLoading, setAuthLoading] = useState(Boolean(api.getToken()));
    const [authError, setAuthError] = useState("");
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
    }, []);

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

    const completeAuth = async (authAction) => {
        setAuthLoading(true);
        setAuthError("");

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
            <section className="landing-screen">
                <div className="landing-content">
                    <h2>AI Inference Gateway</h2>
                    <p>
                        Submit prompts, track queued inference work, and monitor token
                        usage through a PostgreSQL-backed gateway.
                    </p>
                    <button
                        type="button"
                        className="submit-button landing-button"
                        onClick={() => setLandingStarted(true)}
                    >
                        Start
                    </button>
                </div>
            </section>
        );
    }

    if (!currentUser) {
        return (
            <AuthForm
                onLogin={handleLogin}
                onRegister={handleRegister}
                loading={authLoading}
                error={authError}
            />
        );
    }

    if (bootLoading) {
        return (
            <div className="dashboard-layout">
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
        );
    }

    return (
        <div className="dashboard-stack">
            <section className="session-bar">
                <div>
                    <span className="session-label">Signed in</span>
                    <strong>{currentUser.username || authUser?.email || currentUser.email}</strong>
                    <span>{currentUser.email || authUser?.email}</span>
                </div>
                <button type="button" className="logout-button" onClick={handleLogout}>
                    Logout
                </button>
            </section>

            <div className="dashboard-layout">
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
    );
};

export default Dashboard;
