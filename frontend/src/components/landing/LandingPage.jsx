import { useRef, useState } from "react";
import logoMark from "../../assets/logo.png";
import LandingFeatures from "./LandingFeatures";
import LandingHero from "./LandingHero";
import LandingModels from "./LandingModels";
import LandingPricing from "./LandingPricing";

const LandingPage = ({ onStart, onLogin }) => {
    const [landingNavHidden, setLandingNavHidden] = useState(false);
    const landingPageRef = useRef(null);

    const handleLandingScroll = (event) => {
        setLandingNavHidden(event.currentTarget.scrollTop > 24);
    };

    const handleLandingBrandClick = (event) => {
        event.preventDefault();
        setLandingNavHidden(false);
        landingPageRef.current?.scrollTo({ top: 0, behavior: "auto" });
        window.history.replaceState({}, "", window.location.pathname);
    };

    return (
        <div
            className="landing-page"
            onScroll={handleLandingScroll}
            ref={landingPageRef}
        >
            <header
                className={`landing-nav${landingNavHidden ? " landing-nav-hidden" : ""}`}
                aria-label="Primary navigation"
            >
                <a
                    className="landing-brand"
                    href="#top"
                    onClick={handleLandingBrandClick}
                    aria-label="AI Inference Gateway"
                >
                    <img src={logoMark} alt="" className="landing-logo" />
                    <span>AI Inference Gateway</span>
                </a>

                <nav className="landing-menu" aria-label="Page sections">
                    <a href="#features">Features</a>
                    <a href="#models">Models</a>
                    <a href="#pricing">Pricing</a>
                </nav>

                <button
                    type="button"
                    className="landing-login"
                    onClick={onLogin}
                >
                    Login
                </button>
            </header>

            <LandingHero onStart={onStart} />
            <LandingFeatures />
            <LandingModels />
            <LandingPricing onStart={onStart} />
        </div>
    );
};

export default LandingPage;
