import "../styles/components/EmptyState.css";

const EmptyState = ({ title, description }) => {
    return (
        <div className="empty-state">
            <div className="empty-state-copy">
                <p className="empty-state-title">{title}</p>
                {description && (
                    <p className="empty-state-description">{description}</p>
                )}
            </div>
        </div>
    );
};

export default EmptyState;
