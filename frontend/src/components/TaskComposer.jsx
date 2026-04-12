import SectionCard from "./SectionCard";

export default function TaskComposer({
    users,
    models,
    selectedUser,
    selectedModel,
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
}) {
    return (
        <SectionCard as="aside" className="control-panel">
            <form className="control-form" onSubmit={onSubmit}>
                <section className="control-section">
                    <label className="field-label" htmlFor="user-select">
                        User
                    </label>
                    <select
                        id="user-select"
                        value={selectedUser}
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
                        value={selectedModel}
                        onChange={onModelChange}
                        className="field-input"
                    >
                        <option value="">Select model</option>
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
                        placeholder="Enter prompt"
                        disabled={!selectedUser || !selectedModel || submitLoading}
                    />
                </section>

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
                        !selectedUser ||
                        !selectedModel ||
                        !prompt.trim()
                    }
                >
                    {submitLoading ? "Submitting..." : "Submit"}
                </button>
            </form>
        </SectionCard>
    );
}
