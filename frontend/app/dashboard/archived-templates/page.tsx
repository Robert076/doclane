import { redirect } from "next/navigation";
import { Template } from "@/types";
import ArchivedTemplatesSection from "@/components/Pages/ArchivedTemplatesComponents/ArchivedTemplatesSection";
import PageHeader from "@/components/PageHeader/PageHeader";
import { getTemplates } from "@/lib/api/templates";
import { getCurrentUser } from "@/lib/api/users";

export default async function ArchivedTemplatesPage() {
        const userResponse = await getCurrentUser();
        if (!userResponse.success || !userResponse.data) redirect("/login");

        const user = userResponse.data;
        if (user.role !== "admin" && user.department_id === null) {
                redirect("/dashboard/requests");
        }

        const templateResponse = await getTemplates();
        const templates = (templateResponse.data ?? []) as Template[];
        const archivedTemplates = templates.filter((t) => t.is_closed);

        return (
                <div>
                        <PageHeader
                                title="Șabloane arhivate"
                                subtitle="Restaurează și gestionează șabloanele arhivate."
                        />
                        <ArchivedTemplatesSection templates={archivedTemplates} />
                </div>
        );
}
