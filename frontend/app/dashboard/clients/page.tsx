import ClientsSection from "@/components/ClientComponents/ClientsSection/ClientsSection";
import PageHeader from "@/components/PageHeader/PageHeader";
import { notFound } from "next/navigation";
import ClientsActions from "@/components/Pages/ClientsComponents/ClientsActions";
import { getClientsByProfessional } from "@/lib/api/users";

export default async function ClientsPage() {
        const clientsResponse = await getClientsByProfessional();
        if (!clientsResponse.success || !clientsResponse.data) {
                notFound();
        }

        const clients = clientsResponse.data;

        return (
                <div>
                        <PageHeader
                                title="Solicitanți"
                                subtitle="Administrează şi gestionează solicitanții tăi."
                        />
                        <ClientsActions />
                        <ClientsSection clients={clients} />
                </div>
        );
}
