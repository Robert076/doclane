import { User } from "@/types";
import "./ProfessionalView.css";
import getDocumentRequests from "@/lib/getDocumentRequests";
import RequestsSection from "@/components/RequestsSection/RequestsSection";

interface ProfessionalViewProps {
  user: User;
}

export default async function ProfessionalView({ user }: ProfessionalViewProps) {
  const requests = await getDocumentRequests(user.role);

  return (
    <div className="professional-view">
      <header className="professional-header">
        <h1 className="overview-h1">Welcome back, {user.email.split("@")[0]}</h1>
        <p className="overview-p">You have {requests.length} active document requests.</p>
      </header>

      <RequestsSection requests={requests} />
    </div>
  );
}
