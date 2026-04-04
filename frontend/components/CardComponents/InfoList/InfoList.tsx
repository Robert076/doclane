import "./InfoList.css";

interface InfoListProps {
        children: React.ReactNode;
}

export default function InfoList({ children }: InfoListProps) {
        return <div className="info-list">{children}</div>;
}
