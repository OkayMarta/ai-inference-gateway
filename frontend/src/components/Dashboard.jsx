import AuthForm from "./AuthForm";
import EmptyState from "./EmptyState";
import LandingPage from "./landing/LandingPage";
import DashboardDrawer from "./layout/DashboardDrawer";
import MobileDashboardBar from "./layout/MobileDashboardBar";
import SessionBar from "./layout/SessionBar";
import SectionCard from "./SectionCard";
import TaskComposer from "./TaskComposer";
import TaskList from "./TaskList";
import useAuth from "../hooks/useAuth";
import useDashboardData from "../hooks/useDashboardData";
import useMobileDrawer from "../hooks/useMobileDrawer";
import useTasksPolling from "../hooks/useTasksPolling";
import { getInitials } from "../utils/userUtils";

const Dashboard = () => {
    const {
        authError,
        authLoading,
        authSuccess,
        authUser,
        backToLanding,
        currentUser,
        forgotPassword,
        hasAuthToken,
        landingStarted,
        login,
        logout,
        passwordResetToken,
        register,
        resetPassword,
        setCurrentUser,
        startLanding,
    } = useAuth();

    const {
        bootLoading,
        currentModel,
        handleModelChange: updateSelectedModel,
        hasAvailableModels,
        models,
        resetDashboardData,
        screenError,
        selectedModelId,
        setScreenError,
    } = useDashboardData(currentUser);

    const {
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
    } = useTasksPolling({
        currentUser,
        selectedModelId,
        setCurrentUser,
        setScreenError,
    });

    const { closeMenu, mobileMenuOpen, openMenu } = useMobileDrawer();

    const handleLogout = () => {
        logout();
        resetDashboardData();
        resetTasksState();
    };

    const handleModelChange = (event) => {
        updateSelectedModel(event);
        clearSubmitNotices();
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

    if (!landingStarted && !hasAuthToken) {
        return (
            <LandingPage
                onStart={startLanding}
                onLogin={startLanding}
            />
        );
    }

    if (!currentUser) {
        return (
            <AuthForm
                onLogin={login}
                onRegister={register}
                onForgotPassword={forgotPassword}
                onResetPassword={resetPassword}
                onBackToLanding={backToLanding}
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
                    onOpenMenu={openMenu}
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
                        onCloseMenu={closeMenu}
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
