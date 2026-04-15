import PageHeader from "@/components/PageHeader/PageHeader";
import StatsSection from "@/components/Pages/StatsComponents/StatsSection";
import { getStats } from "@/lib/api/stats";
import { getCurrentUser } from "@/lib/api/users";
import { redirect } from "next/navigation";

export default async function StatsPage() {
        const userResponse = await getCurrentUser();
        if (!userResponse.success || !userResponse.data) redirect("/login");
        if (userResponse.data.role !== "admin") redirect("/dashboard/requests");

        const statsResponse = await getStats();
        if (!statsResponse.success || !statsResponse.data) redirect("/dashboard/requests");

        return (
                <div>
                        <PageHeader
                                title="Statistici"
                                subtitle="Vizualizează activitatea platformei în timp real."
                        />
                        <StatsSection stats={statsResponse.data} />
                </div>
        );
}
