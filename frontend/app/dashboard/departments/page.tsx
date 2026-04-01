import DepartmentsSection from "@/components/DepartmentsComponents/DepartmentsSection";
import PageHeader from "@/components/PageHeader/PageHeader";
import { getDepartments } from "@/lib/api/departments";
import { getCurrentUser } from "@/lib/api/users";
import { redirect } from "next/navigation";

export default async function DepartmentsPage() {
        const userResponse = await getCurrentUser();
        if (!userResponse.success || !userResponse.data) redirect("/login");

        const user = userResponse.data;
        if (user.role !== "admin") redirect("/dashboard/requests");

        const departmentsResponse = await getDepartments();
        const departments = departmentsResponse.data ?? [];

        return (
                <div>
                        <PageHeader
                                title="Departamente"
                                subtitle="Administrează departamentele organizației tale."
                        />
                        <DepartmentsSection departments={departments} />
                </div>
        );
}
