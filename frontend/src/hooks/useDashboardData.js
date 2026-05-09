import { useEffect, useMemo, useState } from "react";
import { api } from "../api/client";
import { normalizeList } from "../utils/taskUtils";

const useDashboardData = (currentUser) => {
    const currentUserId = currentUser?.id;
    const [models, setModels] = useState([]);
    const [selectedModelId, setSelectedModelId] = useState("");
    const [bootLoading, setBootLoading] = useState(false);
    const [screenError, setScreenError] = useState("");

    const currentModel = useMemo(
        () => models.find((model) => model.id === selectedModelId) || null,
        [models, selectedModelId],
    );
    const hasAvailableModels = models.length > 0;

    useEffect(() => {
        if (!currentUser) {
            return undefined;
        }

        let active = true;

        const loadBootData = async () => {
            setBootLoading(true);
            setScreenError("");

            try {
                const nextModels = await api.getModels();
                if (!active) {
                    return;
                }
                setModels(normalizeList(nextModels));
            } catch (error) {
                if (active) {
                    setScreenError(error.message);
                }
            } finally {
                if (active) {
                    setBootLoading(false);
                }
            }
        };

        loadBootData();

        return () => {
            active = false;
        };
    }, [currentUserId]);

    const handleModelChange = (event) => {
        setSelectedModelId(event.target.value);
    };

    const resetDashboardData = () => {
        setModels([]);
        setSelectedModelId("");
        setScreenError("");
    };

    return {
        bootLoading,
        currentModel,
        handleModelChange,
        hasAvailableModels,
        models,
        resetDashboardData,
        screenError,
        selectedModelId,
        setScreenError,
    };
};

export default useDashboardData;
