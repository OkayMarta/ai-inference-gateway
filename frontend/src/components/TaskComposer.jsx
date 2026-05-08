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
    const panelTitle = (
        <span className="panel-title-with-icon">
            <svg viewBox="0 0 24 24" aria-hidden="true">
                <path d="M12 3l1.7 4.6L18 9.3l-4.3 1.7L12 16l-1.7-5L6 9.3l4.3-1.7L12 3Z" />
                <path d="M5 15l.9 2.1L8 18l-2.1.9L5 21l-.9-2.1L2 18l2.1-.9L5 15Z" />
                <path d="M19 14l.7 1.6 1.6.7-1.6.7L19 19l-.7-1.6-1.6-.7 1.6-.7L19 14Z" />
            </svg>
            <span>Create task</span>
        </span>
    );

    return (
        <SectionCard as="aside" className="control-panel" title={panelTitle}>
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
                        <svg className="metric-icon metric-icon-balance" viewBox="0 0 24 24" aria-hidden="true">
                            <path d="M4 7h14a3 3 0 0 1 3 3v7a3 3 0 0 1-3 3H4a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h13" />
                            <path d="M18 12h4v4h-4a2 2 0 0 1 0-4Z" />
                            <path d="M6 8h9" />
                        </svg>
                        <div>
                            <span className="metric-label">Balance</span>
                            <span className="metric-value">
                                {currentUser ? currentUser.tokenBalance.toFixed(1) : "-"}
                            </span>
                        </div>
                    </div>
                    <div
                        key={`cost-${metricFlashToken}`}
                        className={`metric-card${balanceAlert ? " metric-card-alert" : ""}`}
                    >
                        <svg className="metric-icon metric-icon-cost" viewBox="0 0 24 24" aria-hidden="true">
                            <ellipse cx="12" cy="6" rx="7" ry="3" />
                            <path d="M5 6v4c0 1.7 3.1 3 7 3s7-1.3 7-3V6" />
                            <path d="M5 10v4c0 1.7 3.1 3 7 3s7-1.3 7-3v-4" />
                            <path d="M5 14v4c0 1.7 3.1 3 7 3s7-1.3 7-3v-4" />
                            <path d="M9 6h6" />
                        </svg>
                        <div>
                            <span className="metric-label">Model cost</span>
                            <span className="metric-value">
                                {currentModel ? currentModel.tokenCost.toFixed(1) : "-"}
                            </span>
                        </div>
                    </div>
                </section>

                <section className="control-section prompt-section">
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
                    <svg viewBox="0 0 24 24" aria-hidden="true">
                        <path d="m22 2-7 20-4-9-9-4 20-7Z" />
                        <path d="M22 2 11 13" />
                    </svg>
                    <span>{submitLoading ? "Submitting..." : "Submit"}</span>
                </button>
            </form>
        </SectionCard>
    );
};

export default TaskComposer;
