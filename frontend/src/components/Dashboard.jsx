import { useEffect, useMemo, useState } from "react";
import { api } from "../api/client";
import logoMark from "../assets/logo.png";
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

const getInitials = (nameOrEmail = "") => {
    const normalized = nameOrEmail.trim();
    if (!normalized) {
        return "AI";
    }

    const nameParts = normalized
        .replace(/@.*$/, "")
        .split(/[\s._-]+/)
        .filter(Boolean);

    return nameParts
        .slice(0, 2)
        .map((part) => part[0]?.toUpperCase())
        .join("");
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
            <div className="landing-page">
                <section className="landing-screen" id="top">
                    <div className="landing-glow landing-glow-primary" aria-hidden="true" />
                    <div className="landing-glow landing-glow-secondary" aria-hidden="true" />

                    <header className="landing-nav" aria-label="Primary navigation">
                        <a className="landing-brand" href="#top" aria-label="AI Inference Gateway">
                            <img src={logoMark} alt="" className="landing-logo" />
                            <span>AI Inference Gateway</span>
                        </a>

                        <nav className="landing-menu" aria-label="Page sections">
                            <a href="#features">Features</a>
                            <a href="#models">Models</a>
                            <a href="#pricing">Pricing</a>
                        </nav>

                        <button
                            type="button"
                            className="landing-login"
                            onClick={() => setLandingStarted(true)}
                        >
                            Login
                        </button>
                    </header>

                    <div className="landing-hero">
                        <div className="landing-copy">
                            <div className="landing-kicker">
                                <svg viewBox="0 0 24 24" aria-hidden="true">
                                    <path d="M13 2 4 14h7l-1 8 10-13h-7l1-7Z" />
                                </svg>
                                <span>AI Request Orchestration & Billing Platform</span>
                            </div>

                            <h1>
                                Orchestrate AI requests with{" "}
                                <span>clarity</span> and <span>control</span>
                            </h1>
                            <p>
                                Submit prompts, choose local models, track queued tasks,
                                and manage token usage through one powerful gateway.
                            </p>

                            <button
                                type="button"
                                className="landing-cta"
                                onClick={() => setLandingStarted(true)}
                            >
                                <span>Get Started</span>
                                <svg viewBox="0 0 24 24" aria-hidden="true">
                                    <path d="M5 12h14" />
                                    <path d="m13 5 7 7-7 7" />
                                </svg>
                            </button>
                        </div>

                        <div className="landing-orbit" aria-hidden="true">
                            <div className="orbit-ring orbit-ring-outer" />
                            <div className="orbit-ring orbit-ring-middle" />
                            <div className="orbit-ring orbit-ring-inner" />
                            <div className="orbit-node orbit-node-blue" />
                            <div className="orbit-node orbit-node-violet" />
                            <div className="orbit-node orbit-node-green" />
                            <div className="orbit-core">
                                <span />
                            </div>

                            <div className="landing-stat landing-stat-balance">
                                <svg viewBox="0 0 24 24">
                                    <ellipse cx="12" cy="5" rx="7" ry="3" />
                                    <path d="M5 5v6c0 1.7 3.1 3 7 3s7-1.3 7-3V5" />
                                    <path d="M5 11v6c0 1.7 3.1 3 7 3s7-1.3 7-3v-6" />
                                </svg>
                                <div>
                                    <span>Token Balance</span>
                                    <strong>1,248,750</strong>
                                    <small>$814.48 USD</small>
                                </div>
                                <div className="landing-progress">
                                    <span />
                                </div>
                            </div>

                            <div className="landing-stat landing-stat-tasks">
                                <svg viewBox="0 0 24 24">
                                    <path d="M9 6h10" />
                                    <path d="M9 12h10" />
                                    <path d="M9 18h10" />
                                    <circle cx="5" cy="6" r="1.5" />
                                    <circle cx="5" cy="12" r="1.5" />
                                    <circle cx="5" cy="18" r="1.5" />
                                </svg>
                                <div>
                                    <span>Queued Tasks</span>
                                    <strong>2,341</strong>
                                    <small>Processing</small>
                                </div>
                            </div>

                            <div className="landing-stat landing-stat-models">
                                <svg viewBox="0 0 24 24">
                                    <path d="m12 2 8 4.5v9L12 20l-8-4.5v-9L12 2Z" />
                                    <path d="M12 11 4.5 6.8" />
                                    <path d="M12 11v8.5" />
                                    <path d="m12 11 7.5-4.2" />
                                </svg>
                                <div>
                                    <span>Local Models</span>
                                    <strong>18</strong>
                                    <small>Available</small>
                                </div>
                            </div>
                        </div>
                    </div>
                </section>

                <section className="landing-section" id="features">
                    <span>Features</span>
                </section>
                <section className="landing-section" id="models">
                    <span>Models</span>
                </section>
                <section className="landing-section" id="pricing">
                    <span>Pricing</span>
                </section>
            </div>
        );
    }

    if (!currentUser) {
        return (
            <AuthForm
                onLogin={handleLogin}
                onRegister={handleRegister}
                onForgotPassword={handleForgotPassword}
                onResetPassword={handleResetPassword}
                onBackToLanding={() => setLandingStarted(false)}
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

    return (
        <div className="dashboard-page" id="top">
            <div className="dashboard-glow dashboard-glow-primary" aria-hidden="true" />
            <div className="dashboard-glow dashboard-glow-secondary" aria-hidden="true" />

            <div className="dashboard-stack">
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

                <section className="session-bar">
                    <div className="session-profile">
                        <span className="session-avatar" aria-hidden="true">
                            {initials}
                        </span>
                        <div className="session-copy">
                            <strong>{displayName}</strong>
                            <span>{displayEmail}</span>
                        </div>
                    </div>
                    <button type="button" className="logout-button" onClick={handleLogout}>
                        <svg viewBox="0 0 24 24" aria-hidden="true">
                            <path d="M10 17l5-5-5-5" />
                            <path d="M15 12H3" />
                            <path d="M14 3h5a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-5" />
                        </svg>
                        <span>Logout</span>
                    </button>
                </section>
            </div>
        </div>
    );
};

export default Dashboard;
