import { useEffect, useMemo, useState } from "react";
import { api } from "../api/client";
import EmptyState from "./EmptyState";
import SectionCard from "./SectionCard";
import TaskComposer from "./TaskComposer";
import TaskList from "./TaskList";
import "../styles/components/Dashboard.css";

function normalizeList(value) {
    return Array.isArray(value) ? value : [];
}

function countTasksByStatus(tasks, status) {
    return tasks.filter((task) => task.status === status).length;
}

export default function Dashboard() {
    const [users, setUsers] = useState([]);
    const [models, setModels] = useState([]);
    const [tasks, setTasks] = useState([]);

    const [selectedUser, setSelectedUser] = useState("");
    const [selectedModel, setSelectedModel] = useState("");
    const [prompt, setPrompt] = useState("");

    const [bootLoading, setBootLoading] = useState(true);
    const [taskLoading, setTaskLoading] = useState(false);
    const [submitLoading, setSubmitLoading] = useState(false);
    const [screenError, setScreenError] = useState("");
    const [submitError, setSubmitError] = useState("");
    const [submitSuccess, setSubmitSuccess] = useState("");

    const currentUser = useMemo(
        () => users.find((user) => user.id === selectedUser) || null,
        [users, selectedUser],
    );

    const currentModel = useMemo(
        () => models.find((model) => model.id === selectedModel) || null,
        [models, selectedModel],
    );

    const sortedTasks = useMemo(
        () =>
            normalizeList(tasks).sort(
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
        if (!selectedUser) {
            setTasks([]);
            setTaskLoading(false);
            return;
        }

        let active = true;

        const refreshUserData = async () => {
            setTaskLoading(true);
            setScreenError("");

            try {
                const [nextTasks, nextUsers] = await Promise.all([
                    api.getTasks(selectedUser),
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
                if (active) {
                    setTaskLoading(false);
                }
            }
        };

        refreshUserData();
        const intervalId = setInterval(refreshUserData, 2000);

        return () => {
            active = false;
            clearInterval(intervalId);
        };
    }, [selectedUser]);

    const handleUserChange = (event) => {
        setSelectedUser(event.target.value);
        setSubmitError("");
        setSubmitSuccess("");
        setScreenError("");
    };

    const handleModelChange = (event) => {
        setSelectedModel(event.target.value);
        setSubmitError("");
        setSubmitSuccess("");
    };

    const handlePromptChange = (event) => {
        setPrompt(event.target.value);
    };

    const handleSubmit = async (event) => {
        event.preventDefault();
        if (!selectedUser || !selectedModel || !prompt.trim()) {
            return;
        }

        setSubmitLoading(true);
        setSubmitError("");
        setSubmitSuccess("");

        try {
            await api.submitTask(selectedUser, selectedModel, prompt.trim());

            const [nextTasks, nextUsers] = await Promise.all([
                api.getTasks(selectedUser),
                api.getUsers(),
            ]);

            setTasks(normalizeList(nextTasks));
            setUsers(normalizeList(nextUsers));
            setPrompt("");
            setSubmitSuccess("Task submitted.");
        } catch (error) {
            setSubmitError(error.message);
        } finally {
            setSubmitLoading(false);
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
                selectedUser={selectedUser}
                selectedModel={selectedModel}
                prompt={prompt}
                currentUser={currentUser}
                currentModel={currentModel}
                screenError={screenError}
                submitError={submitError}
                submitSuccess={submitSuccess}
                submitLoading={submitLoading}
                onUserChange={handleUserChange}
                onModelChange={handleModelChange}
                onPromptChange={handlePromptChange}
                onSubmit={handleSubmit}
            />

            <TaskList
                models={models}
                selectedUser={selectedUser}
                screenError={screenError}
                taskLoading={taskLoading}
                sortedTasks={sortedTasks}
                queuedCount={queuedCount}
                processingCount={processingCount}
                completedCount={completedCount}
            />
        </div>
    );
}
