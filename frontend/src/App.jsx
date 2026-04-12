import Dashboard from "./components/Dashboard";

export default function App() {
    return (
        <div className="min-h-screen relative">
            <div className="fixed inset-0 -z-10 overflow-hidden">
                <div className="absolute -top-40 -right-40 w-96 h-96 bg-indigo-500/20 rounded-full blur-[120px]" />
                <div className="absolute top-1/2 -left-40 w-96 h-96 bg-rose-500/10 rounded-full blur-[120px]" />
            </div>

            <header className="border-b border-white/5 bg-slate-950/50 backdrop-blur-md">
                <div className="max-w-5xl mx-auto px-6 py-5">
                    <h1 className="text-2xl font-bold gradient-text">
                        AI Inference Gateway
                    </h1>
                    <p className="text-white/40 text-sm">
                        Балансувальник навантаження ШІ-моделей
                    </p>
                </div>
            </header>

            <main className="max-w-5xl mx-auto px-6 py-8">
                <Dashboard />
            </main>
        </div>
    );
}
