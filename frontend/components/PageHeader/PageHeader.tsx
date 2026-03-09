import "./PageHeader.css";

interface ProfessionalHeaderProps {
        title: string;
        subtitle: string;
}

const ProfessionalHeader: React.FC<ProfessionalHeaderProps> = ({ title, subtitle }) => {
        return (
                <header className="professional-header">
                        <h1 className="overview-h1">{title}</h1>
                        <p className="overview-p">{subtitle}</p>
                </header>
        );
};

export default ProfessionalHeader;
