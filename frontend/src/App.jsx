import Dashboard from "./components/Dashboard";

export default function App() {
    return (
        <div className="app-shell">
            <header className="app-header">
                <h1 className="app-title">AI Inference Gateway</h1>
            </header>

            <main className="app-main">
                <Dashboard />
            </main>
        </div>
    );
}
