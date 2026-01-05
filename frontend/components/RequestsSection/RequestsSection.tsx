import { DocumentRequest } from "@/types";
import Request from "../Request/Request";
import "./RequestsSection.css";

interface RequestsSectionProps {
  requests: DocumentRequest[];
}

const RequestsSection: React.FC<RequestsSectionProps> = ({ requests }) => {
  if (requests.length === 0) {
    return (
      <div className="requests-empty">
        <p>No document requests found.</p>
      </div>
    );
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
