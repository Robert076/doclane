import getMyClients from "@/lib/getClients";
import ButtonPrimary from "@/components/Buttons/ButtonPrimary/ButtonPrimary";

import "./style.css";
import ClientsSection from "@/components/Client/ClientsSection/ClientsSection";

export default async function ClientsPage() {
  const clients = await getMyClients();

  return (
    <div className="clients-container">
      <header className="clients-header">
        <div>
          <h1 className="overview-h1">My Clients</h1>
          <p className="overview-p">Manage and view your assigned clients.</p>
        </div>
      </header>

      <ClientsSection clients={clients} />
    </div>
  );
}
