import logoMark from "../../assets/logo.png";

const SessionControls = ({ displayEmail, displayName, initials, onLogout }) => (
    <>
        <div className="session-profile">
            <span className="session-avatar" aria-hidden="true">
                {initials}
            </span>
            <div className="session-copy">
                <strong>{displayName}</strong>
                <span>{displayEmail}</span>
            </div>
        </div>
        <button type="button" className="logout-button" onClick={onLogout}>
            <svg viewBox="0 0 24 24" aria-hidden="true">
                <path d="M10 17l5-5-5-5" />
                <path d="M15 12H3" />
                <path d="M14 3h5a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-5" />
            </svg>
            <span>Logout</span>
        </button>
    </>
);

const SessionBar = ({
    displayEmail,
    displayName,
    initials,
    onLogout,
    variant = "desktop",
}) => {
    if (variant === "drawer") {
        return (
            <section className="session-bar session-bar-drawer">
                <SessionControls
                    displayEmail={displayEmail}
                    displayName={displayName}
                    initials={initials}
                    onLogout={onLogout}
                />
            </section>
        );
    }

    return (
        <section className="session-bar session-bar-desktop">
            <div className="session-brand" aria-label="AI Inference Gateway">
                <img src={logoMark} alt="" className="dashboard-logo" />
                <span>AI Inference Gateway</span>
            </div>
            <div className="session-actions">
                <SessionControls
                    displayEmail={displayEmail}
                    displayName={displayName}
                    initials={initials}
                    onLogout={onLogout}
                />
            </div>
        </section>
    );
};

export default SessionBar;
