import Link from "next/link";
import "./AddRequestForClientHeader.css";

interface AddRequestForClientHeaderProps {
  id: string;
}

const AddRequestForClientHeader: React.FC<AddRequestForClientHeaderProps> = ({ id }) => {
  return (
    <div className="add-request-for-client-header">
      <Link href="/dashboard" className="back-link">
        ← Back to clients
      </Link>
      <h1>Add a request for client {id}</h1>
    </div>
  );
};

export default AddRequestForClientHeader;
