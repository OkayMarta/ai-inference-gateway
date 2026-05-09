import { useEffect, useState } from "react";

const useMobileDrawer = () => {
    const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

    useEffect(() => {
        if (!mobileMenuOpen) {
            return undefined;
        }

        const handleKeyDown = (event) => {
            if (event.key === "Escape") {
                setMobileMenuOpen(false);
            }
        };

        window.addEventListener("keydown", handleKeyDown);

        return () => window.removeEventListener("keydown", handleKeyDown);
    }, [mobileMenuOpen]);

    return {
        closeMenu: () => setMobileMenuOpen(false),
        mobileMenuOpen,
        openMenu: () => setMobileMenuOpen(true),
    };
};

export default useMobileDrawer;
