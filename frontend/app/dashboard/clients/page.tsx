import getMyClients from "@/lib/getClients";
import "./style.css";
import InvitationCodeGenerator from "@/components/Invitation/InvitationCodeGenerator/InvitationCodeGenerator";
import InvitationCodesModal from "@/components/Invitation/InvitationCodesModal/InvitationCodesModal";
import ClientsSection from "@/components/ClientComponents/ClientsSection/ClientsSection";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";

export default async function ClientsPage() {
        const clients = await getMyClients();

        return (
                <div className="clients-container">
                        <header className="clients-header">
                                <div>
                                        <h1 className="overview-h1">Clients</h1>
                                        <p className="overview-p">
                                                Manage and view your assigned clients.
                                        </p>
                                </div>
                                <div className="header-actions">
                                        <InvitationCodeGenerator />
                                        <InvitationCodesModal />
                                </div>
                        </header>
                        <ClientsSection clients={clients} />
                </div>
        );
}
