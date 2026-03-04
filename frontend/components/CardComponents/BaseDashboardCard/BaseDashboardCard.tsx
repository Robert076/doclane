import "./BaseDashboardCard.css";

interface BaseDashboardCardProps {
        header?: React.ReactNode;
        title: React.ReactNode;
        children: React.ReactNode;
        footer?: React.ReactNode;
        isHighlighted?: boolean;
}

const BaseDashboardCard: React.FC<BaseDashboardCardProps> = ({
        header,
        title,
        children,
        footer,
        isHighlighted,
}) => {
        return (
                <div className={`dashboard-card ${isHighlighted ? "is-highlighted" : ""}`}>
                        {header && <div className="card-header">{header}</div>}

                        <h3 className="card-title">{title}</h3>

                        <div className="card-body">{children}</div>

                        {footer && <div className="card-footer">{footer}</div>}
                </div>
        );
};

export default BaseDashboardCard;
