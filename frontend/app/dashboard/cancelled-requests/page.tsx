import PageHeader from "@/components/PageHeader/PageHeader";
import CancelledRequestsSection from "@/components/CancelledRequestsComponents/CancelledRequestsSection";
import { getCancelledRequests } from "@/lib/api/requests";
import { getCurrentUser } from "@/lib/api/users";
import { redirect } from "next/navigation";

export default async function CancelledRequestsPage() {
        const userResponse = await getCurrentUser();
        if (!userResponse.success || !userResponse.data) redirect("/login");

        const user = userResponse.data;
        if (user.role !== "admin" && user.department_id === null)
                redirect("/dashboard/requests");

        const res = await getCancelledRequests();
        const requests = res.data ?? [];

        return (
                <div>
                        <PageHeader
                                title="Dosare retrase"
                                subtitle="Vizualizează dosarele retrase de utilizatori."
                        />
                        <CancelledRequestsSection requests={requests} user={user} />
                </div>
        );
}
