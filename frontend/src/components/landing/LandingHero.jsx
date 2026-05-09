const LandingHero = ({ onStart }) => (
    <section className="landing-screen" id="top">
        <div className="landing-glow landing-glow-primary" aria-hidden="true" />
        <div className="landing-glow landing-glow-secondary" aria-hidden="true" />

        <div className="landing-hero">
            <div className="landing-copy">
                <div className="landing-kicker">
                    <svg viewBox="0 0 24 24" aria-hidden="true">
                        <path d="M13 2 4 14h7l-1 8 10-13h-7l1-7Z" />
                    </svg>
                    <span>AI Request Orchestration & Billing Platform</span>
                </div>

                <h1>
                    Orchestrate AI requests with{" "}
                    <span>clarity</span> and <span>control</span>
                </h1>
                <p>
                    Submit prompts, choose local models, track queued tasks,
                    and manage token usage through one powerful gateway.
                </p>

                <button
                    type="button"
                    className="landing-cta"
                    onClick={onStart}
                >
                    <span>Get Started</span>
                    <svg viewBox="0 0 24 24" aria-hidden="true">
                        <path d="M5 12h14" />
                        <path d="m13 5 7 7-7 7" />
                    </svg>
                </button>
            </div>

            <div className="landing-orbit" aria-hidden="true">
                <div className="orbit-ring orbit-ring-outer" />
                <div className="orbit-ring orbit-ring-middle" />
                <div className="orbit-ring orbit-ring-inner" />
                <div className="orbit-node orbit-node-blue" />
                <div className="orbit-node orbit-node-violet" />
                <div className="orbit-node orbit-node-green" />
                <div className="orbit-core">
                    <span />
                </div>

                <div className="landing-stat landing-stat-balance">
                    <svg viewBox="0 0 24 24">
                        <ellipse cx="12" cy="5" rx="7" ry="3" />
                        <path d="M5 5v6c0 1.7 3.1 3 7 3s7-1.3 7-3V5" />
                        <path d="M5 11v6c0 1.7 3.1 3 7 3s7-1.3 7-3v-6" />
                    </svg>
                    <div>
                        <span>Token Balance</span>
                        <strong>1,248,750</strong>
                        <small>$814.48 USD</small>
                    </div>
                    <div className="landing-progress">
                        <span />
                    </div>
                </div>

                <div className="landing-stat landing-stat-tasks">
                    <svg viewBox="0 0 24 24">
                        <path d="M9 6h10" />
                        <path d="M9 12h10" />
                        <path d="M9 18h10" />
                        <circle cx="5" cy="6" r="1.5" />
                        <circle cx="5" cy="12" r="1.5" />
                        <circle cx="5" cy="18" r="1.5" />
                    </svg>
                    <div>
                        <span>Queued Tasks</span>
                        <strong>2,341</strong>
                        <small>Processing</small>
                    </div>
                </div>

                <div className="landing-stat landing-stat-models">
                    <svg viewBox="0 0 24 24">
                        <path d="m12 2 8 4.5v9L12 20l-8-4.5v-9L12 2Z" />
                        <path d="M12 11 4.5 6.8" />
                        <path d="M12 11v8.5" />
                        <path d="m12 11 7.5-4.2" />
                    </svg>
                    <div>
                        <span>Local Models</span>
                        <strong>18</strong>
                        <small>Available</small>
                    </div>
                </div>
            </div>
        </div>
    </section>
);

export default LandingHero;
