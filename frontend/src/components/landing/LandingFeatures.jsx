import { landingFeatureItems } from "./landingData";
import { renderLandingFeatureIcon } from "./landingIcons";

const LandingFeatures = () => (
    <section className="landing-features" id="features">
        <div className="features-orbit" aria-hidden="true">
            <div className="features-orbit-ring" />
            <span className="features-orbit-node features-orbit-node-blue" />
            <span className="features-orbit-node features-orbit-node-violet" />
            <span className="features-orbit-node features-orbit-node-green" />
        </div>

        <div className="features-content">
            <div className="features-heading">
                <h2>
                    Everything you need
                    <span>to run <strong>AI requests</strong></span>
                </h2>
                <p>
                    Orchestrate prompts, manage billing, and monitor every step
                    of your inference pipeline - all in one place.
                </p>
            </div>

            <div className="features-grid">
                {landingFeatureItems.map((feature) => (
                    <article
                        className={`feature-card feature-card-${feature.tone}`}
                        key={feature.title}
                    >
                        <div className="feature-icon">
                            {renderLandingFeatureIcon(feature.icon)}
                        </div>
                        <div className="feature-card-copy">
                            <h3>{feature.title}</h3>
                            <p>{feature.description}</p>
                        </div>
                    </article>
                ))}
            </div>
        </div>
    </section>
);

export default LandingFeatures;
