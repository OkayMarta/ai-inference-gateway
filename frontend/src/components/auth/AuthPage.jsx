import { useEffect, useState } from "react";
import logoMark from "../../assets/logo.png";
import AuthTabs from "./AuthTabs";
import PasswordInput from "./PasswordInput";

const PASSWORD_REQUIREMENTS_MESSAGE =
    "Use 8+ characters with a letter and a number.";

const isValidEmail = (email) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);

const isStrongPassword = (password) => {
    return password.length >= 8 && /[A-Za-z]/.test(password) && /\d/.test(password);
};

const AuthPage = ({
    onLogin,
    onRegister,
    onForgotPassword,
    onResetPassword,
    onBackToLanding,
    onClearMessages,
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
    const [validationError, setValidationError] = useState("");

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
        setValidationError("");

        const trimmedEmail = email.trim();
        const trimmedUsername = username.trim();

        if (isRegister && !trimmedUsername) {
            setValidationError("Username is required.");
            return;
        }

        if (!isReset && !trimmedEmail) {
            setValidationError("Enter your email address.");
            return;
        }

        if (!isReset && !isValidEmail(trimmedEmail)) {
            setValidationError("Enter a valid email address.");
            return;
        }

        if (isForgot) {
            onForgotPassword({ email: trimmedEmail });
            return;
        }

        if (!password) {
            setValidationError("Enter your password.");
            return;
        }

        if ((isRegister || isReset) && !isStrongPassword(password)) {
            setValidationError(PASSWORD_REQUIREMENTS_MESSAGE);
            return;
        }

        if (isReset) {
            onResetPassword({ token: resetToken, newPassword: password });
            return;
        }

        if (isRegister) {
            onRegister({
                username: trimmedUsername,
                email: trimmedEmail,
                password,
            });
            return;
        }

        onLogin({
            email: trimmedEmail,
            password,
        });
    };

    const switchMode = (nextMode) => {
        setMode(nextMode);
        setPassword("");
        setShowPassword(false);
        setValidationError("");
        onClearMessages?.();
    };

    const clearValidationError = () => {
        setValidationError("");
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
    const formError = validationError || error;
    const authCardClassName = `auth-card auth-card-${mode}`;

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

            <div className={authCardClassName}>
                {isForgot || isReset ? (
                    <div className="auth-mode-heading">
                        <h2>{title}</h2>
                    </div>
                ) : (
                    <AuthTabs isRegister={isRegister} onSwitchMode={switchMode} />
                )}

                <form className="auth-form" onSubmit={handleSubmit} noValidate>
                    {isRegister && (
                        <label className="auth-field">
                            <span className="field-label">Username</span>
                            <input
                                className="field-input"
                                placeholder="Choose a username"
                                value={username}
                                onChange={(event) => {
                                    setUsername(event.target.value);
                                    clearValidationError();
                                }}
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
                                onChange={(event) => {
                                    setEmail(event.target.value);
                                    clearValidationError();
                                }}
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
                                onChange={(event) => {
                                    setPassword(event.target.value);
                                    clearValidationError();
                                }}
                                showPassword={showPassword}
                                onToggleVisibility={() =>
                                    setShowPassword((value) => !value)
                                }
                                minLength={isRegister || isReset ? 8 : undefined}
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

                    {formError && <div className="notice notice-error">{formError}</div>}
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
