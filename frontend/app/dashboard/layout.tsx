import Sidebar from "@/components/Sidebar/Sidebar";
import { getCurrentUser } from "@/lib/auth";
import { redirect } from "next/navigation";
import "./dashboard-layout.css";

export default async function DashboardLayout({ children }: { children: React.ReactNode }) {
  const user = await getCurrentUser();

  if (!user) redirect("/login");

  return (
    <div className="dashboard-container">
      <Sidebar user={user} />
      <section className="dashboard-content">{children}</section>
    </div>
  );
}
