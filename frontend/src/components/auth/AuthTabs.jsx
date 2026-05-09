const AuthTabs = ({ isRegister, onSwitchMode }) => (
    <div className="auth-tabs" role="tablist" aria-label="Authentication mode">
        <button
            type="button"
            className={`auth-tab${!isRegister ? " auth-tab-active" : ""}`}
            onClick={() => onSwitchMode("login")}
        >
            Login
        </button>
        <button
            type="button"
            className={`auth-tab${isRegister ? " auth-tab-active" : ""}`}
            onClick={() => onSwitchMode("register")}
        >
            Register
        </button>
    </div>
);

export default AuthTabs;
