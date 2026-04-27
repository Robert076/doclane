"use client";
import { useState } from "react";
import { Tag } from "@/types";
import Modal from "@/components/Modals/Modal";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import { createTag, updateTag, deleteTag } from "@/lib/api/tags";
import toast from "react-hot-toast";
import { useRouter } from "next/navigation";
import "./ManageTagsModal.css";
import Input from "@/components/InputComponents/Input";

interface ManageTagsModalProps {
        isOpen: boolean;
        onClose: () => void;
        tags: Tag[];
}

export default function ManageTagsModal({ isOpen, onClose, tags }: ManageTagsModalProps) {
        const router = useRouter();
        const [name, setName] = useState("");
        const [color, setColor] = useState("#6366f1");
        const [editingTag, setEditingTag] = useState<Tag | null>(null);
        const [isLoading, setIsLoading] = useState(false);

        const resetForm = () => {
                setName("");
                setColor("#6366f1");
                setEditingTag(null);
        };

        const handleSubmit = async () => {
                if (!name.trim()) {
                        toast.error("Numele tagului este obligatoriu.");
                        return;
                }
                setIsLoading(true);
                if (editingTag) {
                        const res = await updateTag(editingTag.id, {
                                name: name.trim(),
                                color,
                        });
                        if (!res.success) {
                                toast.error(res.message ?? "Eroare la actualizare.");
                        } else {
                                toast.success("Tag actualizat!");
                                resetForm();
                                router.refresh();
                        }
                } else {
                        const res = await createTag({ name: name.trim(), color });
                        if (!res.success) {
                                toast.error(res.message ?? "Eroare la creare.");
                        } else {
                                toast.success("Tag creat!");
                                resetForm();
                                router.refresh();
                        }
                }
                setIsLoading(false);
        };

        const handleDelete = async (id: number) => {
                setIsLoading(true);
                const res = await deleteTag(id);
                if (!res.success) {
                        toast.error(res.message ?? "Eroare la ștergere.");
                } else {
                        toast.success("Tag șters!");
                        if (editingTag?.id === id) resetForm();
                        router.refresh();
                }
                setIsLoading(false);
        };

        const handleEdit = (tag: Tag) => {
                setEditingTag(tag);
                setName(tag.name);
                setColor(tag.color);
        };

        return (
                <Modal
                        isOpen={isOpen}
                        onClose={() => {
                                resetForm();
                                onClose();
                        }}
                        title="Gestionează taguri"
                        hideFooter
                >
                        <div className="manage-tags-form">
                                <div className="manage-tags-row">
                                        <Input
                                                placeholder="Nume tag..."
                                                value={name}
                                                onChange={(e) => setName(e.target.value)}
                                                fullWidth
                                        />
                                        <div className="manage-tags-color-picker">
                                                <input
                                                        type="color"
                                                        className="manage-tags-color-input"
                                                        value={color}
                                                        onChange={(e) =>
                                                                setColor(e.target.value)
                                                        }
                                                />
                                                <span
                                                        className="manage-tags-color-preview"
                                                        style={{ backgroundColor: color }}
                                                />
                                        </div>
                                </div>
                                {name.trim() && (
                                        <div className="manage-tags-preview">
                                                <span
                                                        className="manage-tags-badge"
                                                        style={
                                                                {
                                                                        "--tag-color": color,
                                                                } as React.CSSProperties
                                                        }
                                                >
                                                        {name.trim()}
                                                </span>
                                        </div>
                                )}
                                <div className="manage-tags-form-actions">
                                        {editingTag && (
                                                <ButtonPrimary
                                                        text="Anulează"
                                                        variant="ghost"
                                                        fullWidth
                                                        onClick={resetForm}
                                                />
                                        )}
                                        <ButtonPrimary
                                                text={editingTag ? "Salvează" : "Adaugă"}
                                                variant="primary"
                                                fullWidth
                                                onClick={handleSubmit}
                                                disabled={isLoading}
                                        />
                                </div>
                        </div>

                        <div className="manage-tags-list">
                                {tags.length === 0 ? (
                                        <p className="modal-subtext">Nu există taguri.</p>
                                ) : (
                                        tags.map((tag) => (
                                                <div key={tag.id} className="manage-tags-item">
                                                        <span
                                                                className="manage-tags-badge"
                                                                style={
                                                                        {
                                                                                "--tag-color":
                                                                                        tag.color,
                                                                        } as React.CSSProperties
                                                                }
                                                        >
                                                                {tag.name}
                                                        </span>
                                                        <div className="manage-tags-item-actions">
                                                                <button
                                                                        type="button"
                                                                        className="manage-tags-action"
                                                                        onClick={() =>
                                                                                handleEdit(tag)
                                                                        }
                                                                >
                                                                        Editează
                                                                </button>
                                                                <button
                                                                        type="button"
                                                                        className="manage-tags-action manage-tags-action--delete"
                                                                        onClick={() =>
                                                                                handleDelete(
                                                                                        tag.id,
                                                                                )
                                                                        }
                                                                >
                                                                        Șterge
                                                                </button>
                                                        </div>
                                                </div>
                                        ))
                                )}
                        </div>
                </Modal>
        );
}
