import { User } from "@/types";
import RequestsSection from "@/components/RequestsSection/RequestsSection";
import ProfessionalHeader from "./_components/ProfessionalHeader";
import getDocumentRequests from "@/lib/api/getDocumentRequests";

interface ProfessionalViewProps {
  user: User;
}

export default async function ProfessionalView({ user }: ProfessionalViewProps) {
  const requests = await getDocumentRequests(user.role);

  return (
    <div className="professional-view">
      <ProfessionalHeader user={user} length={requests.length} />
      <RequestsSection user={user} requests={requests} />
    </div>
  );
}
