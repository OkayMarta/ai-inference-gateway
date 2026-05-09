import { landingPricingPlans } from "./landingData";

const LandingPricing = ({ onStart }) => (
    <section className="landing-pricing" id="pricing">
        <div className="pricing-content">
            <div className="pricing-heading">
                <h2>
                    Simple <strong>credit-based pricing</strong>
                </h2>
                <p>
                    All requests are paid using platform credits.
                    <span>Choose the plan that fits your usage.</span>
                </p>

                <div className="pricing-rate-card">
                    <div className="pricing-rate-icon">
                        <svg viewBox="0 0 24 24" aria-hidden="true">
                            <path d="M13 2 4 14h7l-1 8 10-13h-7l1-7Z" />
                        </svg>
                    </div>
                    <div>
                        <p>
                            <strong>qwen2.5:1.5b</strong>
                            <span>3 units per request</span>
                        </p>
                        <p>
                            <strong>tinyllama:latest</strong>
                            <span>3 units per request</span>
                        </p>
                    </div>
                </div>
            </div>

            <div className="pricing-grid">
                {landingPricingPlans.map((plan) => (
                    <article
                        className={`pricing-card${plan.recommended ? " pricing-card-recommended" : ""}`}
                        key={plan.name}
                    >
                        {plan.recommended && (
                            <div className="pricing-recommended">
                                <svg viewBox="0 0 24 24" aria-hidden="true">
                                    <path d="m12 2.5 2.8 5.7 6.2.9-4.5 4.4 1.1 6.2-5.6-3-5.6 3 1.1-6.2L3 9.1l6.2-.9L12 2.5Z" />
                                </svg>
                                <span>Recommended</span>
                            </div>
                        )}

                        <h3>{plan.name}</h3>
                        <strong className="pricing-price">{plan.price}</strong>
                        <p className="pricing-units">{plan.units}</p>
                        <p className="pricing-description">{plan.description}</p>
                        <button
                            type="button"
                            className="pricing-button"
                            onClick={onStart}
                        >
                            {plan.cta}
                        </button>
                    </article>
                ))}
            </div>

            <div className="pricing-note">
                <span>
                    <svg viewBox="0 0 24 24" aria-hidden="true">
                        <path d="M12 3 5 6v5c0 4.4 2.9 8.4 7 10 4.1-1.6 7-5.6 7-10V6l-7-3Z" />
                        <path d="m9.5 12 1.8 1.8 3.7-4" />
                    </svg>
                </span>
                <p>Credits are deducted when a task is created.</p>
            </div>
        </div>
    </section>
);

export default LandingPricing;
