export const renderLandingFeatureIcon = (icon) => {
    if (icon === "billing") {
        return (
            <svg viewBox="0 0 24 24" aria-hidden="true">
                <ellipse cx="12" cy="5" rx="7" ry="3" />
                <path d="M5 5v6c0 1.7 3.1 3 7 3s7-1.3 7-3V5" />
                <path d="M5 11v6c0 1.7 3.1 3 7 3s7-1.3 7-3v-6" />
            </svg>
        );
    }

    if (icon === "tasks") {
        return (
            <svg viewBox="0 0 24 24" aria-hidden="true">
                <path d="M9 6h10" />
                <path d="M9 12h10" />
                <path d="M9 18h10" />
                <circle cx="5" cy="6" r="1.4" />
                <circle cx="5" cy="12" r="1.4" />
                <circle cx="5" cy="18" r="1.4" />
            </svg>
        );
    }

    if (icon === "history") {
        return (
            <svg viewBox="0 0 24 24" aria-hidden="true">
                <circle cx="12" cy="12" r="8.5" />
                <path d="M12 7.5V12" />
                <path d="m12 12 3.5 2.1" />
            </svg>
        );
    }

    return (
        <svg viewBox="0 0 24 24" aria-hidden="true">
            <path d="M13 2 4 14h7l-1 8 10-13h-7l1-7Z" />
        </svg>
    );
};

export const renderLandingModelIcon = () => (
    <svg viewBox="0 0 24 24" aria-hidden="true">
        <path d="m12 2 8 4.5v9L12 20l-8-4.5v-9L12 2Z" />
        <path d="M12 11 4.5 6.8" />
        <path d="M12 11v8.5" />
        <path d="m12 11 7.5-4.2" />
        <path d="m7.8 14 4.2-2.4 4.2 2.4" />
    </svg>
);
