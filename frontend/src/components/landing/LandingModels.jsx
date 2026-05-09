import { landingModelItems } from "./landingData";
import { renderLandingModelIcon } from "./landingIcons";

const LandingModels = () => (
    <section className="landing-models" id="models">
        <div className="models-content">
            <div className="models-heading">
                <h2>
                    Available <strong>local models</strong>
                </h2>
                <p>
                    Choose the best local model for your request. You can switch
                    models anytime before sending.
                </p>
            </div>

            <div className="models-grid">
                {landingModelItems.map((model) => (
                    <article
                        className={`model-card${model.recommended ? " model-card-recommended" : ""}`}
                        key={model.name}
                    >
                        {model.recommended && (
                            <div className="model-recommended">
                                <svg viewBox="0 0 24 24" aria-hidden="true">
                                    <path d="m12 2.5 2.8 5.7 6.2.9-4.5 4.4 1.1 6.2-5.6-3-5.6 3 1.1-6.2L3 9.1l6.2-.9L12 2.5Z" />
                                </svg>
                                <span>Recommended</span>
                            </div>
                        )}

                        <div className="model-card-main">
                            <div className="model-icon">
                                {renderLandingModelIcon()}
                            </div>
                            <div className="model-card-copy">
                                <h3>{model.name}</h3>
                                <div className="model-badges">
                                    <span>
                                        <svg viewBox="0 0 24 24" aria-hidden="true">
                                            <path d="M6 5h12" />
                                            <path d="M6 12h12" />
                                            <path d="M6 19h12" />
                                            <rect x="4" y="3" width="16" height="4" rx="1.5" />
                                            <rect x="4" y="10" width="16" height="4" rx="1.5" />
                                            <rect x="4" y="17" width="16" height="4" rx="1.5" />
                                        </svg>
                                        Local model
                                    </span>
                                    <span>
                                        <svg viewBox="0 0 24 24" aria-hidden="true">
                                            <path d="M13 2 4 14h7l-1 8 10-13h-7l1-7Z" />
                                        </svg>
                                        3 units / request
                                    </span>
                                </div>
                                <p>{model.description}</p>
                            </div>
                        </div>

                        <div className="model-best-for">
                            <span>Best for</span>
                            <div>
                                {model.tags.map((tag) => (
                                    <small key={tag}>{tag}</small>
                                ))}
                            </div>
                        </div>
                    </article>
                ))}
            </div>
        </div>
    </section>
);

export default LandingModels;
