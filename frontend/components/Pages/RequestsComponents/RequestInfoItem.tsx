import HighlightText from "@/components/OtherComponents/HighlightText/HighlightText";

interface RequestInfoItemProps {
        label: string;
        value: string;
        searchTerm?: string;
}

export default function RequestInfoItem({ label, value, searchTerm }: RequestInfoItemProps) {
        return (
                <p className="request-info-item">
                        <span className="request-label">{label}</span>
                        <span className="request-value">
                                <HighlightText text={value} search={searchTerm} />
                        </span>
                </p>
        );
}
