import HighlightText from "@/components/OtherComponents/HighlightText/HighlightText";

export default function ClientInfoItem(
        label: string,
        text: string,
        searchTerm: string | undefined,
) {
        return (
                <p className="client-info-item">
                        <span className="client-label">{label}</span>
                        <span className="client-value">
                                <HighlightText text={text} search={searchTerm} />
                        </span>
                </p>
        );
}
