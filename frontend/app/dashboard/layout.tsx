import Sidebar from "@/components/OtherComponents/Sidebar/Sidebar";
import { getCurrentUser } from "@/lib/api/api";
import { redirect } from "next/navigation";
import { UserProvider } from "@/context/UserContext"; // Importă provider-ul
import "./dashboard-layout.css";

export default async function DashboardLayout({ children }: { children: React.ReactNode }) {
        const response = await getCurrentUser();

        if (!response.data) redirect("/login");

        const user = response.data;

        return (
                <UserProvider user={user}>
                        <div className="dashboard-container">
                                <Sidebar />
                                <section className="dashboard-content">{children}</section>
                        </div>
                </UserProvider>
        );
}
