import "./SectionCard.css";

export default function SectionCard({
    as: Component = "section",
    className = "",
    title,
    subtitle,
    rightSlot,
    children,
}) {
    const cardClassName = ["panel", className].filter(Boolean).join(" ");

    return (
        <Component className={cardClassName}>
            {(title || subtitle || rightSlot) && (
                <div className="section-card-header">
                    <div className="section-card-heading">
                        {title && <h2 className="section-card-title">{title}</h2>}
                        {subtitle && (
                            <p className="section-card-subtitle">{subtitle}</p>
                        )}
                    </div>
                    {rightSlot && (
                        <div className="section-card-actions">{rightSlot}</div>
                    )}
                </div>
            )}
            {children}
        </Component>
    );
}
