import TemplatesSection from "../../../components/Pages/TemplatesComponents/TemplatesSection";
import PageHeader from "@/components/PageHeader/PageHeader";
import TemplatesActions from "../../../components/Pages/TemplatesComponents/TemplatesActions";
import { getTemplates } from "@/lib/api/templates";
import { getCurrentUser } from "@/lib/api/users";
import { getDepartments } from "@/lib/api/departments";
import { notFound } from "next/navigation";
import { redirect } from "next/navigation";
import { getTags } from "@/lib/api/tags";

export default async function TemplatesPage() {
        const userResponse = await getCurrentUser();
        if (!userResponse.success || !userResponse.data) redirect("/login");
        const user = userResponse.data;

        const [templatesResponse, departmentsResponse, tagsResponse] = await Promise.all([
                getTemplates(),
                user.role === "admin" ? getDepartments() : Promise.resolve({ data: [] }),
                user.role === "admin" ? getTags() : Promise.resolve({ data: [] }),
        ]);

        if (!templatesResponse.success || !templatesResponse.data) notFound();

        const departments = departmentsResponse.data ?? [];
        const tags = tagsResponse.data ?? [];

        return (
                <div>
                        <PageHeader
                                title="Şabloanele tale"
                                subtitle="Administrează şi gestionează şabloanele tale."
                        />
                        <TemplatesSection
                                templates={templatesResponse.data}
                                isAdmin={user.role === "admin"}
                                userDepartmentId={user.department_id ?? null}
                                departments={departments}
                                tags={tags}
                        />
                </div>
        );
}
