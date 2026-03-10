"use client";
import Modal from "@/components/Modals/Modal";
import React from "react";

interface ConfirmDeleteTemplateModalProps {
        isOpen: boolean;
        onClose: () => void;
        onConfirm: () => void;
        templateTitle?: string;
}

const DeleteTemplateModal: React.FC<ConfirmDeleteTemplateModalProps> = ({
        isOpen,
        onClose,
        onConfirm,
        templateTitle,
}) => {
        return (
                <Modal
                        isOpen={isOpen}
                        onClose={onClose}
                        onConfirm={onConfirm}
                        title="Şterge şablon definitiv"
                >
                        <p>
                                Eşti sigur că vrei să ştergi definitiv şablonul
                                {templateTitle ? (
                                        <>
                                                {" "}
                                                <strong>„{templateTitle}"</strong>
                                        </>
                                ) : (
                                        " acest şablon"
                                )}
                                ? Această acţiune este ireversibilă.
                        </p>
                </Modal>
        );
};

export default DeleteTemplateModal;
