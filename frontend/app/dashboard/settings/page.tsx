import PageHeader from "@/components/PageHeader/PageHeader";
import ProfileSection from "@/components/Pages/SettingsComponents/ProfileSection";
import { getCurrentUser } from "@/lib/api/users";
import { redirect } from "next/navigation";

export default async function SettingsPage() {
        const userResponse = await getCurrentUser();
        if (!userResponse.success || !userResponse.data) redirect("/login");

        return (
                <div>
                        <PageHeader
                                title="Setări"
                                subtitle="Gestionează informațiile contului tău."
                        />
                        <ProfileSection user={userResponse.data} />
                </div>
        );
}
