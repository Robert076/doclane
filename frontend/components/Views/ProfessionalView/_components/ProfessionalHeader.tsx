import "./ProfessionalHeader.css";

interface ProfessionalHeaderProps {
  email: string;
  length: number;
}

const ProfessionalHeader: React.FC<ProfessionalHeaderProps> = ({ email, length }) => {
  return (
    <header className="professional-header">
      <h1 className="overview-h1">Welcome back, {email.split("@")[0]}</h1>
      <p className="overview-p">You have {length} active document requests.</p>
    </header>
  );
};

export default ProfessionalHeader;
