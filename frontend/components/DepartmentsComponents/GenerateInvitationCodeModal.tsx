"use client";
import { useState } from "react";
import Modal from "@/components/Modals/Modal";
import { generateInvitationCode } from "@/lib/api/invitation_codes";
import toast from "react-hot-toast";
import "./GenerateInvitationCodeModal.css";

interface Props {
        isOpen: boolean;
        onClose: () => void;
        departmentId: number;
}

type Step = "confirm" | "generated";

export default function GenerateInvitationCodeModal({ isOpen, onClose, departmentId }: Props) {
        const [step, setStep] = useState<Step>("confirm");
        const [generatedCode, setGeneratedCode] = useState<string | null>(null);

        const handleConfirm = async () => {
                const response = await generateInvitationCode(departmentId, 7);
                if (response.success && response.data) {
                        setGeneratedCode(response.data.code);
                        setStep("generated");
                } else {
                        toast.error(response.message ?? "Eroare la generarea codului.");
                        onClose();
                }
        };

        const handleClose = () => {
                setStep("confirm");
                setGeneratedCode(null);
                onClose();
        };

        const handleCopy = () => {
                if (!generatedCode) return;
                navigator.clipboard.writeText(generatedCode);
                toast.success("Cod copiat!");
        };

        if (step === "confirm") {
                return (
                        <Modal
                                isOpen={isOpen}
                                onClose={handleClose}
                                onConfirm={handleConfirm}
                                title="Adaugă membru"
                                closeOnConfirm={false}
                        >
                                <p>
                                        Dacă continui, se va genera un cod de invitație pentru
                                        acest departament. Codul va expira în{" "}
                                        <strong>7 zile</strong> și poate fi folosit o singură
                                        dată.
                                </p>
                        </Modal>
                );
        }

        return (
                <Modal
                        isOpen={isOpen}
                        onClose={handleClose}
                        onConfirm={handleClose}
                        title="Cod generat cu succes"
                >
                        <div className="generated-code-wrapper">
                                <p className="generated-code-label">
                                        Trimite acest cod persoanei pe care vrei să o adaugi.
                                        Expiră în <strong>7 zile</strong> și poate fi folosit o
                                        singură dată.
                                </p>
                                <div className="generated-code-box">
                                        <span className="generated-code">
                                                {typeof window !== "undefined"
                                                        ? `${window.location.origin}/register/invite?code=${generatedCode}`
                                                        : ""}
                                        </span>
                                        <button
                                                className="copy-btn"
                                                onClick={() => {
                                                        const link = `${window.location.origin}/register/invite?code=${generatedCode}`;
                                                        navigator.clipboard.writeText(link);
                                                        toast.success("Link copiat!");
                                                }}
                                        >
                                                Copiază
                                        </button>
                                </div>
                        </div>
                </Modal>
        );
}
