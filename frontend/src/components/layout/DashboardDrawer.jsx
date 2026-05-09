import logoMark from "../../assets/logo.png";

const DashboardDrawer = ({
    children,
    mobileMenuOpen,
    onCloseMenu,
    sidebarId,
}) => (
    <>
        <button
            type="button"
            className="dashboard-drawer-backdrop"
            onClick={onCloseMenu}
            aria-label="Close dashboard menu"
            tabIndex={mobileMenuOpen ? 0 : -1}
        />

        <div className="dashboard-sidebar" id={sidebarId}>
            <div className="dashboard-drawer-header">
                <a
                    className="dashboard-drawer-brand"
                    href="#top"
                    aria-label="AI Inference Gateway"
                >
                    <img src={logoMark} alt="" className="dashboard-logo" />
                    <span>AI Inference Gateway</span>
                </a>
                <button
                    type="button"
                    className="dashboard-drawer-close"
                    onClick={onCloseMenu}
                    aria-label="Close dashboard menu"
                >
                    <svg viewBox="0 0 24 24" aria-hidden="true">
                        <path d="M18 6 6 18" />
                        <path d="M6 6l12 12" />
                    </svg>
                </button>
            </div>

            {children}
        </div>
    </>
);

export default DashboardDrawer;
