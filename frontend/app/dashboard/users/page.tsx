import PageHeader from "@/components/PageHeader/PageHeader";
import UsersSection from "@/components/Pages/UsersComponents/UsersSection";
import { getCurrentUser, getUsers } from "@/lib/api/users";
import { redirect } from "next/navigation";

export default async function UsersPage() {
        const userResponse = await getCurrentUser();
        if (!userResponse.success || !userResponse.data) redirect("/login");

        const user = userResponse.data;
        if (user.role !== "admin" && user.department_id === null)
                redirect("/dashboard/requests");

        const usersResponse = await getUsers();
        const users = usersResponse.data ?? [];

        return (
                <div>
                        <PageHeader
                                title="Utilizatori"
                                subtitle="Toți utilizatorii înregistrați în platformă."
                        />
                        <UsersSection users={users} />
                </div>
        );
}
