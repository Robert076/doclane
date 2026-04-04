import PageHeader from "@/components/PageHeader/PageHeader";
import ArchivedRequestsSection from "@/components/Pages/ArchivedRequestsComponents/ArchivedRequestsSection";
import { getAllRequests, getRequestsByDepartment } from "@/lib/api/requests";
import { getCurrentUser } from "@/lib/api/users";
import { redirect } from "next/navigation";
import { Request } from "@/types";

export default async function ArchivedRequestsPage() {
        const userResponse = await getCurrentUser();
        if (!userResponse.success || !userResponse.data) redirect("/login");

        const user = userResponse.data;

        if (user.role !== "admin" && user.department_id === null) {
                redirect("/dashboard/requests");
        }

        let requests: Request[] = [];

        if (user.role === "admin") {
                const res = await getAllRequests();
                requests = res.data ?? [];
        } else {
                const res = await getRequestsByDepartment(user.department_id!);
                requests = res.data ?? [];
        }

        const archivedRequests = requests.filter((r) => r.is_closed);

        return (
                <div>
                        <PageHeader
                                title="Dosare arhivate"
                                subtitle="Vizualizează și gestionează dosarele arhivate."
                        />
                        <ArchivedRequestsSection requests={archivedRequests} user={user} />
                </div>
        );
}
