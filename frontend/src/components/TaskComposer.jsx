import SectionCard from "./SectionCard";

const TaskComposer = ({
    users,
    models,
    hasAvailableModels,
    selectedUserId,
    selectedModelId,
    prompt,
    currentUser,
    currentModel,
    screenError,
    submitError,
    submitSuccess,
    submitLoading,
    onUserChange,
    onModelChange,
    onPromptChange,
    onSubmit,
}) => {
    const noModelsNotice =
        "Немає доступних моделей. Переконайтеся, що Ollama запущена і має завантажені моделі.";

    return (
        <SectionCard as="aside" className="control-panel">
            <form className="control-form" onSubmit={onSubmit}>
                <section className="control-section">
                    <label className="field-label" htmlFor="user-select">
                        User
                    </label>
                    <select
                        id="user-select"
                        value={selectedUserId}
                        onChange={onUserChange}
                        className="field-input"
                    >
                        <option value="">Select user</option>
                        {users.map((user) => (
                            <option key={user.id} value={user.id}>
                                {user.username}
                            </option>
                        ))}
                    </select>
                </section>

                <section className="control-section">
                    <label className="field-label" htmlFor="model-select">
                        Model
                    </label>
                    <select
                        id="model-select"
                        value={selectedModelId}
                        onChange={onModelChange}
                        className="field-input"
                        disabled={!hasAvailableModels}
                    >
                        <option value="">
                            {hasAvailableModels
                                ? "Select model"
                                : "No models available"}
                        </option>
                        {models.map((model) => (
                            <option key={model.id} value={model.id}>
                                {model.name}
                            </option>
                        ))}
                    </select>
                </section>

                <section className="metrics-grid">
                    <div className="metric-card">
                        <span className="metric-label">Balance</span>
                        <span className="metric-value">
                            {currentUser ? currentUser.tokenBalance.toFixed(1) : "-"}
                        </span>
                    </div>
                    <div className="metric-card">
                        <span className="metric-label">Model cost</span>
                        <span className="metric-value">
                            {currentModel ? currentModel.tokenCost.toFixed(1) : "-"}
                        </span>
                    </div>
                </section>

                <section className="control-section">
                    <label className="field-label" htmlFor="prompt-input">
                        Prompt
                    </label>
                    <textarea
                        id="prompt-input"
                        value={prompt}
                        onChange={onPromptChange}
                        className="field-input field-textarea"
                        placeholder={
                            hasAvailableModels
                                ? "Enter prompt"
                                : "Task submission is unavailable without models"
                        }
                        disabled={
                            !hasAvailableModels ||
                            !selectedUserId ||
                            !selectedModelId ||
                            submitLoading
                        }
                    />
                </section>

                {!hasAvailableModels && (
                    <div className="notice notice-warning">{noModelsNotice}</div>
                )}
                {screenError && <div className="notice notice-error">{screenError}</div>}
                {submitError && <div className="notice notice-error">{submitError}</div>}
                {submitSuccess && (
                    <div className="notice notice-success">{submitSuccess}</div>
                )}

                <button
                    type="submit"
                    className="submit-button"
                    disabled={
                        submitLoading ||
                        !hasAvailableModels ||
                        !selectedUserId ||
                        !selectedModelId ||
                        !prompt.trim()
                    }
                >
                    {submitLoading ? "Submitting..." : "Submit"}
                </button>
            </form>
        </SectionCard>
    );
};

export default TaskComposer;
