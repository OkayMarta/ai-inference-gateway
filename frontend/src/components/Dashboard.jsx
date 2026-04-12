import { useState, useEffect, useRef } from "react";
import { api } from "../api/client";

export default function Dashboard() {
    const [users, setUsers] = useState([]);
    const [models, setModels] = useState([]);
    const [tasks, setTasks] = useState([]);

    const [selectedUser, setSelectedUser] = useState("");
    const [selectedModel, setSelectedModel] = useState("");
    const [prompt, setPrompt] = useState("");

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    const messagesEndRef = useRef(null);

    // 1. При першому завантаженні отримуємо списки юзерів і моделей
    useEffect(() => {
        api.getUsers().then(setUsers).catch(console.error);
        api.getModels().then(setModels).catch(console.error);
    }, []);

    // 2. Опитування (Polling) - кожні 2 секунди оновлюємо задачі обраного юзера
    useEffect(() => {
        if (!selectedUser) {
            setTasks([]);
            return;
        }

        const fetchTasks = () => {
            api.getTasks(selectedUser).then(setTasks).catch(console.error);
            // Також оновлюємо баланс юзера у списку
            api.getUsers().then(setUsers).catch(console.error);
        };

        fetchTasks(); // Викликаємо одразу
        const interval = setInterval(fetchTasks, 2000); // І потім кожні 2 сек

        return () => clearInterval(interval); // Очищаємо при зміні юзера
    }, [selectedUser]);

    // Автоскрол до останнього повідомлення
    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }, [tasks]);

    // Відправка форми
    const handleSubmit = async (e) => {
        e.preventDefault();
        if (!selectedUser || !selectedModel || !prompt.trim()) return;

        setLoading(true);
        setError("");
        try {
            await api.submitTask(selectedUser, selectedModel, prompt.trim());
            setPrompt(""); // Очищаємо поле
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    const currentUser = users.find((u) => u.id === selectedUser);
    const selectedModelInfo = models.find((m) => m.id === selectedModel);

    // Сортуємо задачі від старих до нових для чату
    const sortedTasks = [...tasks].sort(
        (a, b) => new Date(a.createdAt) - new Date(b.createdAt),
    );

    return (
        <div className="space-y-6">
            {/* Верхня панель: Вибір юзера і баланс */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="glass-card p-5">
                    <label className="block text-white/50 text-sm font-medium mb-2">
                        Обрати користувача
                    </label>
                    <select
                        value={selectedUser}
                        onChange={(e) => setSelectedUser(e.target.value)}
                        className="input-field"
                    >
                        <option value="">Оберіть...</option>
                        {users.map((u) => (
                            <option key={u.id} value={u.id}>
                                {u.username}
                            </option>
                        ))}
                    </select>
                </div>

                <div className="glass-card p-5 flex items-center justify-between">
                    <div>
                        <p className="text-white/50 text-sm font-medium">
                            Баланс токенів
                        </p>
                        <p className="text-3xl font-bold gradient-text mt-1">
                            {currentUser
                                ? currentUser.tokenBalance.toFixed(1)
                                : "—"}
                        </p>
                    </div>
                    <div className="text-4xl">💎</div>
                </div>
            </div>

            {/* Вибір моделі */}
            <div>
                <h2 className="text-sm font-semibold text-white/50 mb-3">
                    Доступні ШІ-моделі
                </h2>
                <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
                    {models.map((m) => (
                        <button
                            key={m.id}
                            onClick={() => setSelectedModel(m.id)}
                            className={`glass-card p-4 text-left transition-all duration-300 ${
                                selectedModel === m.id
                                    ? "border-indigo-500 shadow-lg shadow-indigo-500/20 bg-indigo-500/10"
                                    : "hover:border-white/30"
                            }`}
                        >
                            <h3 className="font-bold text-white text-sm">
                                {m.name}
                            </h3>
                            <p className="text-indigo-300 font-semibold text-xs mt-2">
                                {m.tokenCost} токенів
                            </p>
                        </button>
                    ))}
                </div>
            </div>

            {/* Вікно чату (Історія задач) */}
            <div className="glass-card p-4 h-[400px] overflow-y-auto space-y-4">
                {sortedTasks.length === 0 ? (
                    <div className="text-center text-white/30 mt-20">
                        Завдань ще немає. Надішліть перший запит!
                    </div>
                ) : (
                    sortedTasks.map((task) => {
                        const m = models.find((x) => x.id === task.modelId);
                        return (
                            <div key={task.id} className="space-y-2">
                                {/* Запит юзера */}
                                <div className="flex justify-end">
                                    <div className="max-w-[80%] bg-indigo-600/40 border border-indigo-500/30 rounded-2xl rounded-tr-sm px-4 py-2 text-sm">
                                        {task.payload}
                                    </div>
                                </div>
                                {/* Відповідь ШІ */}
                                <div className="flex justify-start">
                                    <div className="max-w-[85%] bg-white/5 border border-white/10 rounded-2xl rounded-tl-sm px-4 py-3 space-y-2">
                                        <div className="flex items-center gap-2">
                                            <span className="text-xs font-bold text-rose-400">
                                                🤖 {m?.name || task.modelId}
                                            </span>
                                            <span
                                                className={`badge ${
                                                    task.status === "Processing"
                                                        ? "badge-processing"
                                                        : task.status ===
                                                            "Completed"
                                                          ? "badge-completed"
                                                          : task.status ===
                                                              "Failed"
                                                            ? "badge-failed"
                                                            : "badge-queued"
                                                }`}
                                            >
                                                {task.status}
                                            </span>
                                        </div>
                                        <p className="text-white/80 text-sm whitespace-pre-wrap">
                                            {task.status === "Completed" ||
                                            task.status === "Failed" ? (
                                                task.result
                                            ) : (
                                                <span className="text-white/30 italic">
                                                    Очікуємо генерацію...
                                                </span>
                                            )}
                                        </p>
                                    </div>
                                </div>
                            </div>
                        );
                    })
                )}
                <div ref={messagesEndRef} />
            </div>

            {/* Форма вводу */}
            <form onSubmit={handleSubmit} className="glass-card p-4">
                {error && (
                    <div className="text-rose-400 text-sm mb-2">{error}</div>
                )}
                <div className="flex gap-3">
                    <input
                        type="text"
                        value={prompt}
                        onChange={(e) => setPrompt(e.target.value)}
                        placeholder="Введіть ваш запит (промпт)..."
                        className="input-field flex-1"
                        disabled={!selectedUser || !selectedModel}
                    />
                    <button
                        type="submit"
                        disabled={
                            loading ||
                            !selectedUser ||
                            !selectedModel ||
                            !prompt.trim()
                        }
                        className="btn-primary"
                    >
                        {loading ? "⏳" : "🚀"} Відправити
                    </button>
                </div>
            </form>
        </div>
    );
}
