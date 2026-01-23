import { User } from "@/types";
import getDocumentRequests from "@/lib/getDocumentRequests";
import RequestsSection from "@/components/RequestsSection/RequestsSection";
import ProfessionalHeader from "./_components/ProfessionalHeader";

interface ProfessionalViewProps {
  user: User;
}

export default async function ProfessionalView({ user }: ProfessionalViewProps) {
  const requests = await getDocumentRequests(user.role);

  return (
    <div className="professional-view">
      <ProfessionalHeader email={user.email} length={requests.length} />
      <RequestsSection requests={requests} />
    </div>
  );
}
