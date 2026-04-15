import { useEffect, useMemo, useState } from "react";
import { api } from "../api/client";
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

const Dashboard = () => {
    const [users, setUsers] = useState([]);
    const [models, setModels] = useState([]);
    const [tasks, setTasks] = useState([]);

    const [selectedUserId, setSelectedUserId] = useState("");
    const [selectedModelId, setSelectedModelId] = useState("");
    const [prompt, setPrompt] = useState("");
    const [statusFilter, setStatusFilter] = useState("");

    const [bootLoading, setBootLoading] = useState(true);
    const [taskLoading, setTaskLoading] = useState(false);
    const [submitLoading, setSubmitLoading] = useState(false);
    const [cancelLoadingTaskId, setCancelLoadingTaskId] = useState("");
    const [screenError, setScreenError] = useState("");
    const [submitError, setSubmitError] = useState("");
    const [submitSuccess, setSubmitSuccess] = useState("");
    const [balanceAlert, setBalanceAlert] = useState("");
    const [metricFlashToken, setMetricFlashToken] = useState(0);

    const currentUser = useMemo(
        () => users.find((user) => user.id === selectedUserId) || null,
        [users, selectedUserId],
    );

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

    const hasBootData = users.length > 0 || models.length > 0;
    const composerScreenError = !hasBootData ? screenError : "";
    const taskScreenError = selectedUserId ? screenError : "";

    useEffect(() => {
        let active = true;

        const loadBootData = async () => {
            setBootLoading(true);
            setScreenError("");

            try {
                const [nextUsers, nextModels] = await Promise.all([
                    api.getUsers(),
                    api.getModels(),
                ]);

                if (!active) {
                    return;
                }

                setUsers(normalizeList(nextUsers));
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
    }, []);

    useEffect(() => {
        if (!selectedUserId) {
            setTasks([]);
            setTaskLoading(false);
            return;
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
                        userId: selectedUserId,
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
                setUsers(normalizeList(nextUsers));
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
    }, [selectedUserId, statusFilter]);

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

    const handleUserChange = (event) => {
        setSelectedUserId(event.target.value);
        setSubmitError("");
        setSubmitSuccess("");
        setBalanceAlert("");
        setScreenError("");
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
        if (!selectedUserId || !selectedModelId || !prompt.trim()) {
            return;
        }

        setSubmitLoading(true);
        setSubmitError("");
        setSubmitSuccess("");
        setBalanceAlert("");

        try {
            await api.submitTask(selectedUserId, selectedModelId, prompt.trim());

            const [nextTasks, nextUsers] = await Promise.all([
                api.getTasks({
                    userId: selectedUserId,
                    status: statusFilter,
                    limit: 20,
                    offset: 0,
                    sort: "created_at_desc",
                }),
                api.getUsers(),
            ]);

            setTasks(normalizeList(nextTasks));
            setUsers(normalizeList(nextUsers));
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
                    userId: selectedUserId,
                    status: statusFilter,
                    limit: 20,
                    offset: 0,
                    sort: "created_at_desc",
                }),
                api.getUsers(),
            ]);

            setTasks(normalizeList(nextTasks));
            setUsers(normalizeList(nextUsers));
        } catch (error) {
            console.error(error);
            window.alert(`Error: ${error.message}`);
        } finally {
            setCancelLoadingTaskId("");
        }
    };

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
        <div className="dashboard-layout">
            <TaskComposer
                users={users}
                models={models}
                hasAvailableModels={hasAvailableModels}
                selectedUserId={selectedUserId}
                selectedModelId={selectedModelId}
                prompt={prompt}
                currentUser={currentUser}
                currentModel={currentModel}
                screenError={composerScreenError}
                submitError={submitError}
                submitSuccess={submitSuccess}
                balanceAlert={balanceAlert}
                metricFlashToken={metricFlashToken}
                submitLoading={submitLoading}
                onDismissSubmitError={() => setSubmitError("")}
                onDismissSubmitSuccess={() => setSubmitSuccess("")}
                onDismissBalanceAlert={() => setBalanceAlert("")}
                onUserChange={handleUserChange}
                onModelChange={handleModelChange}
                onPromptChange={handlePromptChange}
                onSubmit={handleSubmit}
            />

            <TaskList
                models={models}
                selectedUserId={selectedUserId}
                screenError={taskScreenError}
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
    );
};

export default Dashboard;
