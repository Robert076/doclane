import { DocumentRequest } from "@/types";
import Request from "../Request/Request";
import "./RequestsSection.css";
import EmptyRequestsSection from "./_components/EmptyRequestsSection";

interface RequestsSectionProps {
  requests: DocumentRequest[];
}

const RequestsSection: React.FC<RequestsSectionProps> = ({ requests }) => {
  if (requests.length === 0) {
    return <EmptyRequestsSection />;
  }

  return (
    <div className="requests-grid">
      {requests.map((req) => (
        <Request key={req.id} request={req} />
      ))}
    </div>
  );
};

export default RequestsSection;
