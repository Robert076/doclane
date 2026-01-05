import Sidebar from "@/components/Sidebar/Sidebar";
import { getCurrentUser } from "@/lib/auth";
import { redirect } from "next/navigation";
import { UserProvider } from "@/context/UserContext"; // Importă provider-ul
import "./dashboard-layout.css";

export default async function DashboardLayout({ children }: { children: React.ReactNode }) {
  const user = await getCurrentUser();

  if (!user) redirect("/login");

  return (
    <UserProvider user={user}>
      <div className="dashboard-container">
        <Sidebar />
        <section className="dashboard-content">{children}</section>
      </div>
    </UserProvider>
  );
}
