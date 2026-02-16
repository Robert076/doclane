import { Suspense } from "react";
import ClientView from "@/components/ViewComponents/ClientView/ClientView";
import ProfessionalView from "@/components/ViewComponents/ProfessionalView/ProfessionalView";
import { getCurrentUser } from "@/lib/api/api";
import { redirect } from "next/navigation";
import DeactivatedAccountView from "@/components/ViewComponents/DeactivatedAccountView/DeactivatedAccountView";

export const metadata = {
        title: "Dashboard | Doclane",
};

export default async function DashboardPage() {
        const response = await getCurrentUser();

        if (!response.data) {
                redirect("/login");
        }

        const user = response.data;

        if (!user.is_active) {
                return <DeactivatedAccountView />;
        }

        return (
                <main className="min-h-screen bg-gray-50">
                        {(() => {
                                switch (user.role) {
                                        case "PROFESSIONAL":
                                                return <ProfessionalView user={user} />;
                                        case "CLIENT":
                                                return <ClientView user={user} />;
                                        default:
                                                return <p>Unknown role.</p>;
                                }
                        })()}
                </main>
        );
}
