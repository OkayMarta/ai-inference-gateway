import SectionCard from "./SectionCard";

const TaskComposer = ({
    models,
    hasAvailableModels,
    selectedModelId,
    prompt,
    currentUser,
    currentModel,
    screenError,
    submitError,
    submitSuccess,
    balanceAlert,
    metricFlashToken,
    submitLoading,
    onDismissSubmitError,
    onDismissSubmitSuccess,
    onDismissBalanceAlert,
    onModelChange,
    onPromptChange,
    onSubmit,
}) => {
    const noModelsNotice =
        "No models are available. Make sure Ollama is running and has models loaded.";

    return (
        <SectionCard as="aside" className="control-panel">
            <form className="control-form" onSubmit={onSubmit}>
                {(balanceAlert || submitError || submitSuccess) && (
                    <div className="floating-notice-stack" aria-live="polite">
                        {balanceAlert && (
                            <div className="floating-notice floating-notice-error" role="alert">
                                <span>{balanceAlert}</span>
                                <button
                                    type="button"
                                    className="floating-notice-close"
                                    onClick={onDismissBalanceAlert}
                                    aria-label="Dismiss balance alert"
                                >
                                    x
                                </button>
                            </div>
                        )}
                        {submitError && (
                            <div className="floating-notice floating-notice-error" role="alert">
                                <span>{submitError}</span>
                                <button
                                    type="button"
                                    className="floating-notice-close"
                                    onClick={onDismissSubmitError}
                                    aria-label="Dismiss submit error"
                                >
                                    x
                                </button>
                            </div>
                        )}
                        {submitSuccess && (
                            <div
                                className="floating-notice floating-notice-success"
                                role="status"
                            >
                                <span>{submitSuccess}</span>
                                <button
                                    type="button"
                                    className="floating-notice-close"
                                    onClick={onDismissSubmitSuccess}
                                    aria-label="Dismiss submit success"
                                >
                                    x
                                </button>
                            </div>
                        )}
                    </div>
                )}

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
                    <div
                        key={`balance-${metricFlashToken}`}
                        className={`metric-card${balanceAlert ? " metric-card-alert" : ""}`}
                    >
                        <span className="metric-label">Balance</span>
                        <span className="metric-value">
                            {currentUser ? currentUser.tokenBalance.toFixed(1) : "-"}
                        </span>
                    </div>
                    <div
                        key={`cost-${metricFlashToken}`}
                        className={`metric-card${balanceAlert ? " metric-card-alert" : ""}`}
                    >
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
                            !selectedModelId ||
                            submitLoading
                        }
                    />
                </section>

                {!hasAvailableModels && (
                    <div className="notice notice-warning">{noModelsNotice}</div>
                )}
                {screenError && <div className="notice notice-error">{screenError}</div>}

                <button
                    type="submit"
                    className="submit-button"
                    disabled={
                        submitLoading ||
                        !hasAvailableModels ||
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
