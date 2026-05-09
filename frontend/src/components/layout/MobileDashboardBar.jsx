import logoMark from "../../assets/logo.png";

const MobileDashboardBar = ({ mobileMenuOpen, onOpenMenu, sidebarId }) => (
    <header className="dashboard-mobile-bar">
        <button
            type="button"
            className="dashboard-menu-button"
            onClick={onOpenMenu}
            aria-label="Open dashboard menu"
            aria-controls={sidebarId}
            aria-expanded={mobileMenuOpen}
        >
            <svg viewBox="0 0 24 24" aria-hidden="true">
                <path d="M4 7h16" />
                <path d="M4 12h16" />
                <path d="M4 17h16" />
            </svg>
        </button>
        <a className="dashboard-mobile-brand" href="#top" aria-label="AI Inference Gateway">
            <img src={logoMark} alt="" className="dashboard-logo" />
            <span>AI Inference Gateway</span>
        </a>
    </header>
);

export default MobileDashboardBar;
