"use client";
import { useState } from "react";
import Modal from "@/components/Modals/Modal";
import Input from "@/components/InputComponents/Input";
import toast from "react-hot-toast";
import { createDepartment } from "@/lib/api/departments";

interface CreateDepartmentModalProps {
        isOpen: boolean;
        onClose: () => void;
        onCreated: () => void;
}

export default function CreateDepartmentModal({
        isOpen,
        onClose,
        onCreated,
}: CreateDepartmentModalProps) {
        const [name, setName] = useState("");

        const handleConfirm = async () => {
                if (!name.trim()) {
                        toast.error("Numele departamentului este obligatoriu.");
                        return;
                }

                const response = await createDepartment(name.trim());
                if (response.success) {
                        toast.success("Departament creat cu succes!");
                        setName("");
                        onClose();
                        onCreated();
                } else {
                        toast.error(response.message);
                }
        };

        return (
                <Modal
                        isOpen={isOpen}
                        onClose={onClose}
                        onConfirm={handleConfirm}
                        title="Departament nou"
                >
                        <Input
                                label="Nume departament"
                                placeholder="Ex: Administrație publică, Taxe și impozite..."
                                value={name}
                                onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                                        setName(e.target.value)
                                }
                        />
                </Modal>
        );
}
