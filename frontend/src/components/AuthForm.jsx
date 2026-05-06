import { useState } from "react";

const AuthForm = ({ onLogin, onRegister, loading, error }) => {
    const [mode, setMode] = useState("login");
    const [username, setUsername] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [showPassword, setShowPassword] = useState(false);

    const isRegister = mode === "register";

    const handleSubmit = (event) => {
        event.preventDefault();
        if (isRegister) {
            onRegister({
                username: username.trim(),
                email: email.trim(),
                password,
            });
            return;
        }

        onLogin({
            email: email.trim(),
            password,
        });
    };

    const switchMode = (nextMode) => {
        setMode(nextMode);
        setPassword("");
        setShowPassword(false);
    };

    return (
        <section className="auth-shell" aria-label="Authentication">
            <div className="auth-card">
                <div className="auth-tabs" role="tablist" aria-label="Authentication mode">
                    <button
                        type="button"
                        className={`auth-tab${!isRegister ? " auth-tab-active" : ""}`}
                        onClick={() => switchMode("login")}
                    >
                        Login
                    </button>
                    <button
                        type="button"
                        className={`auth-tab${isRegister ? " auth-tab-active" : ""}`}
                        onClick={() => switchMode("register")}
                    >
                        Register
                    </button>
                </div>

                <form className="auth-form" onSubmit={handleSubmit}>
                    {isRegister && (
                        <label className="auth-field">
                            <span className="field-label">Username</span>
                            <input
                                className="field-input"
                                value={username}
                                onChange={(event) => setUsername(event.target.value)}
                                autoComplete="username"
                                required
                            />
                        </label>
                    )}

                    <label className="auth-field">
                        <span className="field-label">Email</span>
                        <input
                            className="field-input"
                            type="email"
                            value={email}
                            onChange={(event) => setEmail(event.target.value)}
                            autoComplete="email"
                            required
                        />
                    </label>

                    <label className="auth-field">
                        <span className="field-label">Password</span>
                        <span className="password-input-wrap">
                            <input
                                className="field-input password-input"
                                type={showPassword ? "text" : "password"}
                                value={password}
                                onChange={(event) => setPassword(event.target.value)}
                                autoComplete={isRegister ? "new-password" : "current-password"}
                                minLength={6}
                                required
                            />
                            <button
                                type="button"
                                className="password-toggle"
                                onClick={() => setShowPassword((value) => !value)}
                                aria-label={showPassword ? "Hide password" : "Show password"}
                                title={showPassword ? "Hide password" : "Show password"}
                            >
                                <svg
                                    className="password-toggle-icon"
                                    viewBox="0 0 24 24"
                                    aria-hidden="true"
                                >
                                    {showPassword ? (
                                        <>
                                            <path d="M3 3l18 18" />
                                            <path d="M10.6 10.6a2 2 0 0 0 2.8 2.8" />
                                            <path d="M9.9 4.2A10.7 10.7 0 0 1 12 4c5 0 8.8 3.4 10 8a11.8 11.8 0 0 1-2.2 4" />
                                            <path d="M6.6 6.6A11.8 11.8 0 0 0 2 12c1.2 4.6 5 8 10 8 1.3 0 2.5-.2 3.6-.7" />
                                        </>
                                    ) : (
                                        <>
                                            <path d="M2 12s3.6-7 10-7 10 7 10 7-3.6 7-10 7S2 12 2 12z" />
                                            <circle cx="12" cy="12" r="3" />
                                        </>
                                    )}
                                </svg>
                            </button>
                        </span>
                    </label>

                    {error && <div className="notice notice-error">{error}</div>}

                    <button type="submit" className="submit-button" disabled={loading}>
                        {loading ? "Please wait..." : isRegister ? "Create account" : "Login"}
                    </button>
                </form>
            </div>
        </section>
    );
};

export default AuthForm;
