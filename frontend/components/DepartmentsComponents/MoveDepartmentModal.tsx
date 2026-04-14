"use client";
import { useState } from "react";
import { Department } from "@/types";
import Modal from "@/components/Modals/Modal";
import SearchBar from "@/components/OtherComponents/SearchBar/SearchBar";
import HighlightText from "@/components/OtherComponents/HighlightText/HighlightText";
import { updateUserDepartment } from "@/lib/api/users";
import { useRouter } from "next/navigation";
import toast from "react-hot-toast";
import "./MoveDepartmentModal.css";

interface Props {
        isOpen: boolean;
        onClose: () => void;
        userId: number;
        userName: string;
        currentDepartmentId: number;
        departments: Department[];
}

export default function MoveDepartmentModal({
        isOpen,
        onClose,
        userId,
        userName,
        currentDepartmentId,
        departments,
}: Props) {
        const router = useRouter();
        const [search, setSearch] = useState("");
        const [selectedId, setSelectedId] = useState<number | null>(null);
        const [isSubmitting, setIsSubmitting] = useState(false);

        const filtered = departments.filter(
                (d) =>
                        d.id !== currentDepartmentId &&
                        d.name.toLowerCase().includes(search.toLowerCase()),
        );

        const handleConfirm = async () => {
                if (!selectedId) {
                        toast.error("Selectează un departament.");
                        return;
                }
                setIsSubmitting(true);
                const response = await updateUserDepartment(userId, selectedId);
                setIsSubmitting(false);
                if (response.success) {
                        toast.success(`${userName} a fost mutat cu succes.`);
                        router.refresh();
                        onClose();
                } else {
                        toast.error(response.message ?? "Eroare la mutarea utilizatorului.");
                }
        };

        const handleClose = () => {
                setSearch("");
                setSelectedId(null);
                onClose();
        };

        return (
                <Modal
                        isOpen={isOpen}
                        onClose={handleClose}
                        onConfirm={handleConfirm}
                        title={`Mută ${userName} în alt departament`}
                        closeOnConfirm={false}
                >
                        <div className="move-department-modal">
                                <SearchBar
                                        value={search}
                                        onChange={setSearch}
                                        placeholder="Caută departament..."
                                        fullWidth
                                />
                                <div className="move-department-list">
                                        {filtered.length === 0 ? (
                                                <p className="move-department-empty">
                                                        Niciun departament găsit.
                                                </p>
                                        ) : (
                                                filtered.map((dept) => (
                                                        <div
                                                                key={dept.id}
                                                                className={`move-department-item ${selectedId === dept.id ? "move-department-item--selected" : ""}`}
                                                                onClick={() =>
                                                                        setSelectedId(dept.id)
                                                                }
                                                        >
                                                                <HighlightText
                                                                        text={dept.name}
                                                                        search={search}
                                                                />
                                                        </div>
                                                ))
                                        )}
                                </div>
                        </div>
                </Modal>
        );
}
