import HighlightText from "@/components/OtherComponents/HighlightText/HighlightText";

export default function ClientInfoItem(
        label: string,
        text: string,
        searchTerm: string | undefined,
        needsHighlighting: boolean,
) {
        return (
                <p className="client-info-item">
                        <span className="client-label">{label}</span>

                        {needsHighlighting && (
                                <span className="client-value">
                                        {" "}
                                        <HighlightText text={text} search={searchTerm} />
                                </span>
                        )}

                        {needsHighlighting === false && (
                                <span className="client-value">{text}</span>
                        )}
                </p>
        );
}
