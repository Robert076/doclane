import { User } from "@/types";
import "./ProfessionalHeader.css";

interface ProfessionalHeaderProps {
  user: User;
  length: number;
}

const ProfessionalHeader: React.FC<ProfessionalHeaderProps> = ({ user, length }) => {
  return (
    <header className="professional-header">
      <h1 className="overview-h1">
        Welcome back, {user.first_name} {user.last_name}
      </h1>
      <p className="overview-p">You have {length} active document requests.</p>
    </header>
  );
};

export default ProfessionalHeader;
