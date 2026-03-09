import Link from "next/link";
import "./AddRequestForClientHeader.css";
import { UI_TEXT } from "@/locales/ro";

interface AddRequestForClientHeaderProps {
        id: string;
}

const AddRequestForClientHeader: React.FC<AddRequestForClientHeaderProps> = ({ id }) => {
        return (
                <div className="add-request-for-client-header">
                        <Link href="/dashboard" className="back-link">
                                {UI_TEXT.common.back}
                        </Link>
                        <h1>Dosar nou</h1>
                        <p>Introdu detaliile noului dosar.</p>
                </div>
        );
};

export default AddRequestForClientHeader;
