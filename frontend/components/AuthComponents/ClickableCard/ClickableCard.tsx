import "./ClickableCard.css";

interface ClickableCardProps {
        text: string;
        icon?: React.ReactNode;
        onClick?: () => void;
}

const ClickableCard: React.FC<ClickableCardProps> = ({ text, icon, onClick }) => {
        return (
                <div className="clickable-card" onClick={onClick}>
                        {icon && <span className="clickable-card-icon">{icon}</span>}
                        <p className="clickable-card-text">{text}</p>
                </div>
        );
};

export default ClickableCard;
