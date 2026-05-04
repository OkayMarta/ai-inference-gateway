import { useState } from "react";

const AuthForm = ({ onLogin, onRegister, loading, error }) => {
    const [mode, setMode] = useState("login");
    const [username, setUsername] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

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
                        <input
                            className="field-input"
                            type="password"
                            value={password}
                            onChange={(event) => setPassword(event.target.value)}
                            autoComplete={isRegister ? "new-password" : "current-password"}
                            minLength={6}
                            required
                        />
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
