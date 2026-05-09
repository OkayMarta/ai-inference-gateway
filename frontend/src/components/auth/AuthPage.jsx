import { useEffect, useState } from "react";
import logoMark from "../../assets/logo.png";
import AuthTabs from "./AuthTabs";
import PasswordInput from "./PasswordInput";

const AuthPage = ({
    onLogin,
    onRegister,
    onForgotPassword,
    onResetPassword,
    onBackToLanding,
    loading,
    error,
    success,
    resetToken,
}) => {
    const [mode, setMode] = useState(resetToken ? "reset" : "login");
    const [username, setUsername] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [showPassword, setShowPassword] = useState(false);

    const isRegister = mode === "register";
    const isForgot = mode === "forgot";
    const isReset = mode === "reset";

    useEffect(() => {
        if (resetToken) {
            setMode("reset");
            setPassword("");
            setShowPassword(false);
        }
    }, [resetToken]);

    const handleSubmit = (event) => {
        event.preventDefault();
        if (isForgot) {
            onForgotPassword({ email: email.trim() });
            return;
        }

        if (isReset) {
            onResetPassword({ token: resetToken, newPassword: password });
            return;
        }

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

    const title = isForgot
        ? "Reset access"
        : isReset
          ? "Set new password"
          : isRegister
            ? "Create account"
            : "Login";

    const passwordLabel = isReset ? "New password" : "Password";
    const passwordPlaceholder = isReset
        ? "Enter a new password"
        : isRegister
          ? "Create a strong password"
          : "Enter your password";

    return (
        <section className="auth-shell" aria-label="Authentication">
            <div className="auth-glow auth-glow-primary" aria-hidden="true" />
            <div className="auth-glow auth-glow-secondary" aria-hidden="true" />

            <button
                type="button"
                className="auth-brand"
                onClick={onBackToLanding}
                aria-label="Back to landing page"
            >
                <img src={logoMark} alt="" className="auth-logo" />
                <span>AI Inference Gateway</span>
            </button>

            <div className="auth-card">
                {isForgot || isReset ? (
                    <div className="auth-mode-heading">
                        <h2>{title}</h2>
                    </div>
                ) : (
                    <AuthTabs isRegister={isRegister} onSwitchMode={switchMode} />
                )}

                <form className="auth-form" onSubmit={handleSubmit}>
                    {isRegister && (
                        <label className="auth-field">
                            <span className="field-label">Username</span>
                            <input
                                className="field-input"
                                placeholder="Choose a username"
                                value={username}
                                onChange={(event) => setUsername(event.target.value)}
                                autoComplete="username"
                                required
                            />
                        </label>
                    )}

                    {!isReset && (
                        <label className="auth-field">
                        <span className="field-label">Email</span>
                        <input
                            className="field-input"
                                type="email"
                                placeholder={
                                    isForgot
                                        ? "Email linked to your account"
                                        : isRegister
                                          ? "you@example.com"
                                          : "Enter your email"
                                }
                            value={email}
                            onChange={(event) => setEmail(event.target.value)}
                            autoComplete="email"
                            required
                        />
                    </label>
                    )}

                    {!isForgot && (
                        <label className="auth-field">
                            <span className="field-label">{passwordLabel}</span>
                        <PasswordInput
                            autoComplete={
                                isRegister || isReset
                                    ? "new-password"
                                    : "current-password"
                            }
                            placeholder={passwordPlaceholder}
                            value={password}
                            onChange={(event) => setPassword(event.target.value)}
                            showPassword={showPassword}
                            onToggleVisibility={() =>
                                setShowPassword((value) => !value)
                            }
                        />
                    </label>
                    )}

                    {!isRegister && !isForgot && !isReset && (
                        <button
                            type="button"
                            className="forgot-password-button"
                            onClick={() => switchMode("forgot")}
                        >
                            Forgot password?
                        </button>
                    )}

                    {error && <div className="notice notice-error">{error}</div>}
                    {success && <div className="notice notice-success">{success}</div>}

                    <button type="submit" className="submit-button" disabled={loading}>
                        {loading
                            ? "Please wait..."
                            : isForgot
                              ? "Send reset link"
                              : isReset
                                ? "Reset password"
                                : isRegister
                                  ? "Create account"
                                  : "Login"}
                    </button>

                    {(isRegister || isForgot || isReset) && (
                        <p className="auth-switch-copy">
                            {isRegister ? "Already have an account? " : ""}
                            <button type="button" onClick={() => switchMode("login")}>
                                Back to login
                            </button>
                        </p>
                    )}
                </form>
            </div>
        </section>
    );
};

export default AuthPage;
