import { User } from "@/types";

interface ProfessionalViewProps {
  user: User;
}

const ProfessionalView: React.FC<ProfessionalViewProps> = ({ user }) => {
  return (
    <div>
      <h1>Welcome back, {user.email.split("@")[0]}</h1>
    </div>
  );
};

export default ProfessionalView;
