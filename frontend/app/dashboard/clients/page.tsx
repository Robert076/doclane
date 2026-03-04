import getMyClients from "@/lib/getClients";
import "./style.css";
import InvitationCodeGenerator from "@/components/ClientComponents/Invitation/InvitationCodeGenerator/InvitationCodeGenerator";
import InvitationCodesModal from "@/components/ClientComponents/Invitation/InvitationCodesModal/InvitationCodesModal";
import ClientsSection from "@/components/ClientComponents/ClientsSection/ClientsSection";
import { UI_TEXT } from "@/locales/ro";

export default async function ClientsPage() {
        const clients = await getMyClients();

        return (
                <div className="clients-container">
                        <header className="clients-header">
                                <div>
                                        <h1 className="overview-h1">
                                                {UI_TEXT.dashboard.professional.headerClients}
                                        </h1>
                                        <p className="overview-p">
                                                {
                                                        UI_TEXT.dashboard.professional
                                                                .subheaderClients
                                                }
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
