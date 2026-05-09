export const getInitials = (nameOrEmail = "") => {
    const normalized = nameOrEmail.trim();
    if (!normalized) {
        return "AI";
    }

    const nameParts = normalized
        .replace(/@.*$/, "")
        .split(/[\s._-]+/)
        .filter(Boolean);

    return nameParts
        .slice(0, 2)
        .map((part) => part[0]?.toUpperCase())
        .join("");
};

export const sameUserSnapshot = (left, right) => {
    if (!left || !right) {
        return false;
    }

    return (
        left.id === right.id &&
        left.username === right.username &&
        left.email === right.email &&
        left.role === right.role &&
        left.tokenBalance === right.tokenBalance
    );
};
