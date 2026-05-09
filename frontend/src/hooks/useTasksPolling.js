import { useCallback, useEffect, useMemo, useState } from "react";
import { api } from "../api/client";
import { countTasksByStatus, normalizeList } from "../utils/taskUtils";
import { sameUserSnapshot } from "../utils/userUtils";
import useAutoDismiss from "./useAutoDismiss";

const TASK_POLL_INTERVAL = 2000;

const useTasksPolling = ({
    currentUser,
    selectedModelId,
    setCurrentUser,
    setScreenError,
}) => {
    const currentUserId = currentUser?.id;
    const [tasks, setTasks] = useState([]);
    const [prompt, setPrompt] = useState("");
    const [statusFilter, setStatusFilter] = useState("");
    const [taskLoading, setTaskLoading] = useState(false);
    const [submitLoading, setSubmitLoading] = useState(false);
    const [cancelLoadingTaskId, setCancelLoadingTaskId] = useState("");
    const [submitError, setSubmitError] = useState("");
    const [submitSuccess, setSubmitSuccess] = useState("");
    const [balanceAlert, setBalanceAlert] = useState("");
    const [metricFlashToken, setMetricFlashToken] = useState(0);

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

    const clearSubmitNotices = useCallback(() => {
        setSubmitError("");
        setSubmitSuccess("");
        setBalanceAlert("");
    }, []);

    useAutoDismiss(submitError, () => setSubmitError(""), 4200);
    useAutoDismiss(submitSuccess, () => setSubmitSuccess(""), 2400);
    useAutoDismiss(balanceAlert, () => setBalanceAlert(""), 4200);

    const applyTasksAndUser = useCallback(
        (nextTasks, nextUsers) => {
            setTasks(normalizeList(nextTasks));
            const refreshedUser = normalizeList(nextUsers).find(
                (user) => user.id === currentUserId,
            );
            if (refreshedUser) {
                setCurrentUser((previousUser) =>
                    sameUserSnapshot(previousUser, refreshedUser)
                        ? previousUser
                        : refreshedUser,
                );
            }
        },
        [currentUserId, setCurrentUser],
    );

    const loadTasksAndUser = useCallback(
        () =>
            Promise.all([
                api.getTasks({
                    userId: currentUserId,
                    status: statusFilter,
                    limit: 20,
                    offset: 0,
                    sort: "created_at_desc",
                }),
                api.getUsers(),
            ]),
        [currentUserId, statusFilter],
    );

    useEffect(() => {
        if (!currentUserId) {
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
                const [nextTasks, nextUsers] = await loadTasksAndUser();

                if (!active) {
                    return;
                }

                applyTasksAndUser(nextTasks, nextUsers);
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
        }, TASK_POLL_INTERVAL);

        return () => {
            active = false;
            clearInterval(intervalId);
        };
    }, [applyTasksAndUser, currentUserId, loadTasksAndUser, setScreenError]);

    const handlePromptChange = (event) => {
        clearSubmitNotices();
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
        clearSubmitNotices();

        try {
            await api.submitTask(selectedModelId, prompt.trim());

            const [nextTasks, nextUsers] = await loadTasksAndUser();

            applyTasksAndUser(nextTasks, nextUsers);
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

            const [nextTasks, nextUsers] = await loadTasksAndUser();

            applyTasksAndUser(nextTasks, nextUsers);
        } catch (error) {
            console.error(error);
            window.alert(`Error: ${error.message}`);
        } finally {
            setCancelLoadingTaskId("");
        }
    };

    const resetTasksState = () => {
        setTasks([]);
        setPrompt("");
        setStatusFilter("");
    };

    return {
        balanceAlert,
        cancelLoadingTaskId,
        cancelledCount,
        clearSubmitNotices,
        completedCount,
        failedCount,
        handleCancelTask,
        handlePromptChange,
        handleStatusFilterChange,
        handleSubmit,
        metricFlashToken,
        processingCount,
        prompt,
        queuedCount,
        resetTasksState,
        setBalanceAlert,
        setSubmitError,
        setSubmitSuccess,
        sortedTasks,
        statusFilter,
        submitError,
        submitLoading,
        submitSuccess,
        taskLoading,
    };
};

export default useTasksPolling;
