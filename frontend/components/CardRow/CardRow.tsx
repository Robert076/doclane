interface CardRowProps {
        label: string | React.ReactNode;
        value: string | React.ReactNode;
        isDescription?: boolean;
}

export default function CardRow({ label, value, isDescription }: CardRowProps) {
        return (
                <div className={isDescription ? "info-description" : "info-row"}>
                        <strong>{label}</strong>
                        {isDescription ? <p>{value}</p> : <span>{value}</span>}
                </div>
        );
}
