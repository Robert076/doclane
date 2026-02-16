import HighlightText from "@/components/HighlightText/HighlightText";

export default function RequestInfoItem(
        searchTerm: string | undefined,
        label: string,
        text: string,
) {
        return (
                <p className="request-info-item">
                        <span className="request-label">{label}</span>
                        <span className="request-value">
                                <HighlightText text={text} search={searchTerm} />
                        </span>
                </p>
        );
}
