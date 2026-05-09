import { useEffect, useRef } from "react";

const useAutoDismiss = (value, onDismiss, delay) => {
    const onDismissRef = useRef(onDismiss);

    useEffect(() => {
        onDismissRef.current = onDismiss;
    }, [onDismiss]);

    useEffect(() => {
        if (!value) {
            return undefined;
        }

        const timeoutId = window.setTimeout(() => {
            onDismissRef.current();
        }, delay);

        return () => window.clearTimeout(timeoutId);
    }, [delay, value]);
};

export default useAutoDismiss;
