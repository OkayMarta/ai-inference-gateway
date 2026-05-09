const PasswordInput = ({
    autoComplete,
    onChange,
    onToggleVisibility,
    placeholder,
    showPassword,
    value,
    minLength,
}) => (
    <span className="password-input-wrap">
        <input
            className="field-input password-input"
            type={showPassword ? "text" : "password"}
            placeholder={placeholder}
            value={value}
            onChange={onChange}
            autoComplete={autoComplete}
            minLength={minLength}
            required
        />
        <button
            type="button"
            className="password-toggle"
            onClick={onToggleVisibility}
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
);

export default PasswordInput;
