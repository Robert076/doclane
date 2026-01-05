import { Suspense } from "react";
import ClientView from "@/components/Views/ClientView/ClientView";
import ProfessionalView from "@/components/Views/ProfessionalView/ProfessionalView";
import { getCurrentUser } from "@/lib/auth";
import { redirect } from "next/navigation";
import DeactivatedAccountView from "@/components/Views/DeactivatedAccountView/DeactivatedAccountView";
import LoadingSkeleton from "@/components/Views/LoadingSkeleton/LoadingSkeleton";

export const metadata = {
  title: "Dashboard | Doclane",
};

export default async function DashboardPage() {
  const user = await getCurrentUser();

  if (!user) {
    redirect("/login");
  }

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
            return <ClientView />;
          default:
            return <p>Unknown role.</p>;
        }
      })()}
    </main>
  );
}
