import PageHeader from "@/components/PageHeader/PageHeader";
import RequestsSection from "@/components/Pages/RequestsComponents/RequestsSection";
import {
        getAllRequests,
        getRequestsByAssignee,
        getRequestsByDepartment,
} from "@/lib/api/requests";
import { getCurrentUser } from "@/lib/api/users";
import { redirect } from "next/navigation";
import { Request } from "@/types";

export default async function RequestsPage() {
        const userResponse = await getCurrentUser();
        if (!userResponse.success || !userResponse.data) redirect("/login");

        const user = userResponse.data;
        let requests: Request[] = [];

        if (user.role === "admin") {
                const res = await getAllRequests();
                requests = res.data ?? [];
        } else if (user.department_id) {
                const res = await getRequestsByDepartment(user.department_id);
                requests = res.data ?? [];
        } else {
                const res = await getRequestsByAssignee(user.id);
                requests = res.data ?? [];
        }

        return (
                <div>
                        <PageHeader
                                title="Dosarele tale"
                                subtitle="Administrează şi gestionează dosarele tale deschise."
                        />
                        <RequestsSection user={user} requests={requests} />
                </div>
        );
}
