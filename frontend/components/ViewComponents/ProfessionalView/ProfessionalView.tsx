import { User } from "@/types";
import RequestsSection from "@/components/RequestComponents/RequestsSection/RequestsSection";
import ProfessionalHeader from "./_components/ProfessionalHeader";
import { getDocumentRequests } from "@/lib/api/api";

interface ProfessionalViewProps {
        user: User;
}

export default async function ProfessionalView({ user }: ProfessionalViewProps) {
        const response = await getDocumentRequests(user.role);

        return (
                <div className="professional-view">
                        <ProfessionalHeader user={user} length={response.data.length} />
                        <RequestsSection user={user} requests={response.data} />
                </div>
        );
}
