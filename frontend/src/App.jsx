import Dashboard from "./components/Dashboard";
import "./styles/app.css";

const App = () => {
    return (
        <div className="app-shell">
            <header className="app-header">
                <div className="app-header-line" aria-hidden="true" />
                <h1 className="app-title">AI Inference Gateway</h1>
                <div className="app-header-line" aria-hidden="true" />
            </header>

            <main className="app-main" aria-label="Dashboard workspace">
                <Dashboard />
            </main>
        </div>
    );
};

export default App;
