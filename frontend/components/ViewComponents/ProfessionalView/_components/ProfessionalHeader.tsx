import { User } from "@/types";
import "./ProfessionalHeader.css";
import { UI_TEXT } from "@/locales/ro";

interface ProfessionalHeaderProps {
        user: User;
        length: number;
}

const ProfessionalHeader: React.FC<ProfessionalHeaderProps> = ({ user, length }) => {
        return (
                <header className="professional-header">
                        <h1 className="overview-h1">
                                {UI_TEXT.dashboard.professional.headerDocumentRequests(
                                        `${user.first_name} ${user.last_name}`,
                                )}
                        </h1>
                        <p className="overview-p">
                                {UI_TEXT.dashboard.professional.subheaderDocumentRequests(
                                        length,
                                )}
                        </p>
                </header>
        );
};

export default ProfessionalHeader;
