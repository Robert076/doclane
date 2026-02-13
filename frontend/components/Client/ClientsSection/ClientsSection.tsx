"use client";
import { useEffect, useState } from "react";
import { User } from "@/types";
import ClientCard from "../ClientCard/ClientCard";
import NotFound from "@/components/NotFound/NotFound";
import "./ClientsSection.css";
import PaginationFooter from "./_components/PaginationFooter";
import SearchBar from "@/components/SearchBar/SearchBar";

interface ClientsSectionProps {
  clients: User[];
}

const ITEMS_PER_PAGE = 12;

const ClientsSection: React.FC<ClientsSectionProps> = ({ clients }) => {
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [searchInput, setSearchInput] = useState<string>("");

  const filteredClients = clients.filter((client) => {
    if (!searchInput) return true;
    const searchLower = searchInput.toLowerCase().trim();

    const fullName = `${client.first_name || ""} ${client.last_name || ""}`.toLowerCase();

    return (
      client.first_name?.toLowerCase().includes(searchLower) ||
      client.last_name?.toLowerCase().includes(searchLower) ||
      fullName.includes(searchLower) ||
      client.email?.toLowerCase().includes(searchLower)
    );
  });

  useEffect(() => {
    setCurrentPage(1);
  }, [searchInput]);

  const totalPages = Math.ceil(clients.length / ITEMS_PER_PAGE);
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
  const endIndex = startIndex + ITEMS_PER_PAGE;

  return (
    <div className="clients-section">
      {filteredClients.length > 0 && (
        <SearchBar
          value={searchInput}
          onChange={setSearchInput}
          placeholder="Search clients..."
        />
      )}
      {filteredClients.length === 0 && (
        <NotFound
          text="No clients found."
          subtext="Start by adding your first client."
          background="#fff"
        />
      )}
      <div className="clients-grid">
        {filteredClients.length > 0 &&
          filteredClients.map((client) => (
            <ClientCard key={client.id} client={client} searchTerm={searchInput} />
          ))}
      </div>

      {totalPages > 1 && (
        <PaginationFooter
          currentPage={currentPage}
          totalPages={totalPages}
          setCurrentPage={setCurrentPage}
        />
      )}
    </div>
  );
};

export default ClientsSection;
