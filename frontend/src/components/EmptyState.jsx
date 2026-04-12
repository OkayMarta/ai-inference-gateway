import "./EmptyState.css";

export default function EmptyState({ title, description }) {
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
}
