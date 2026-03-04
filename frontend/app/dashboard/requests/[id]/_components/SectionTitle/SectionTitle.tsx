import "./SectionTitle.css";

interface SectionTitleProps {
        text: string;
}

const SectionTitle: React.FC<SectionTitleProps> = ({ text }) => {
        return (
                <h2 className="section-title" title={text}>
                        {text}
                </h2>
        );
};

export default SectionTitle;
