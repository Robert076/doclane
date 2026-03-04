"use client";
import { useState, useEffect } from "react";
import { User } from "@/types";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import HighlightText from "@/components/OtherComponents/HighlightText/HighlightText";
import PaginationFooter from "@/components/FileSectionComponents/FileSection/_components/PaginationFooter";
import "./AssignClientModal.css";

interface AssignClientModalProps {
        isOpen: boolean;
        onClose: () => void;
        onConfirm: (client: User) => void;
        clients: User[];
}

const ITEMS_PER_PAGE = 5;

const AssignClientModal: React.FC<AssignClientModalProps> = ({
        isOpen,
        onClose,
        onConfirm,
        clients,
}) => {
        const [search, setSearch] = useState("");
        const [selected, setSelected] = useState<User | null>(null);
        const [currentPage, setCurrentPage] = useState(1);

        useEffect(() => {
                setCurrentPage(1);
        }, [search]);

        if (!isOpen) return null;

        const filtered = clients.filter((client) => {
                const q = search.toLowerCase().trim();
                if (!q) return true;
                const fullName = `${client.first_name} ${client.last_name}`.toLowerCase();
                return (
                        client.first_name.toLowerCase().includes(q) ||
                        client.last_name.toLowerCase().includes(q) ||
                        fullName.includes(q) ||
                        client.email.toLowerCase().includes(q)
                );
        });

        const totalPages = Math.ceil(filtered.length / ITEMS_PER_PAGE);
        const currentClients = filtered.slice(
                (currentPage - 1) * ITEMS_PER_PAGE,
                currentPage * ITEMS_PER_PAGE,
        );

        const handleConfirm = () => {
                if (!selected) return;
                onConfirm(selected);
                onClose();
                setSelected(null);
                setSearch("");
        };

        const handleClose = () => {
                onClose();
                setSelected(null);
                setSearch("");
        };

        return (
                <div className="assign-client-modal-overlay" onClick={handleClose}>
                        <div
                                className="assign-client-modal-content"
                                onClick={(e) => e.stopPropagation()}
                        >
                                <div className="assign-client-modal-header">
                                        <h3>Selectează solicitantul</h3>
                                        <button
                                                className="assign-client-modal-close"
                                                onClick={handleClose}
                                        >
                                                ×
                                        </button>
                                </div>
                                <div className="assign-client-modal-body">
                                        <input
                                                className="assign-client-search"
                                                placeholder="Caută după nume sau email..."
                                                value={search}
                                                onChange={(e) => setSearch(e.target.value)}
                                                autoFocus
                                        />
                                        <div className="assign-client-list">
                                                {filtered.length === 0 && (
                                                        <p className="assign-client-empty">
                                                                Niciun solicitant găsit.
                                                        </p>
                                                )}
                                                {currentClients.map((client) => (
                                                        <div
                                                                key={client.id}
                                                                className={`assign-client-item ${selected?.id === client.id ? "selected" : ""}`}
                                                                onClick={() =>
                                                                        setSelected(client)
                                                                }
                                                        >
                                                                <span className="assign-client-name">
                                                                        <HighlightText
                                                                                text={`${client.first_name} ${client.last_name}`}
                                                                                search={search}
                                                                        />
                                                                </span>
                                                                <span className="assign-client-email">
                                                                        <HighlightText
                                                                                text={
                                                                                        client.email
                                                                                }
                                                                                search={search}
                                                                        />
                                                                </span>
                                                        </div>
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
                                <div className="assign-client-modal-footer">
                                        <ButtonPrimary
                                                text="Anulează"
                                                variant="ghost"
                                                onClick={handleClose}
                                        />
                                        <ButtonPrimary
                                                text="Continuă"
                                                variant="primary"
                                                onClick={handleConfirm}
                                                disabled={!selected}
                                        />
                                </div>
                        </div>
                </div>
        );
};

export default AssignClientModal;
