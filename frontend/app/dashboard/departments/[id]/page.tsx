import PageHeader from "@/components/PageHeader/PageHeader";
import { getDepartments } from "@/lib/api/departments";
import { getUsersByDepartment, getCurrentUser } from "@/lib/api/users";
import { redirect } from "next/navigation";
import DepartmentMembersSection from "@/components/DepartmentsComponents/DepartmentMembersSection";

export default async function DepartmentMembersPage({
        params,
}: {
        params: Promise<{ id: string }>;
}) {
        const { id } = await params;
        const deptId = Number(id);

        const userResponse = await getCurrentUser();
        if (!userResponse.success || !userResponse.data) redirect("/login");
        if (userResponse.data.role !== "admin") redirect("/dashboard/requests");

        const [departmentsResponse, membersResponse] = await Promise.all([
                getDepartments(),
                getUsersByDepartment(deptId),
        ]);

        const department = departmentsResponse.data?.find((d) => d.id === deptId);
        if (!department) redirect("/dashboard/departments");

        const members = membersResponse.data ?? [];

        return (
                <div>
                        <PageHeader
                                title={department.name}
                                subtitle="Membrii acestui departament."
                        />
                        <DepartmentMembersSection
                                members={members}
                                departmentId={deptId}
                                departments={departmentsResponse.data ?? []}
                        />
                </div>
        );
}
