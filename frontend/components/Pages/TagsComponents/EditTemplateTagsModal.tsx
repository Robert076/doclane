// EditTemplateTagsModal.tsx
"use client";
import { useEffect, useState } from "react";
import { Tag, Template } from "@/types";
import Modal from "@/components/Modals/Modal";
import { getTags, setTemplateTags } from "@/lib/api/tags";
import toast from "react-hot-toast";
import { useRouter } from "next/navigation";
import "./EditTemplateTagsModal.css";

const MAX_TAGS = 3;

interface EditTemplateTagsModalProps {
        isOpen: boolean;
        onClose: () => void;
        template: Template;
}

export default function EditTemplateTagsModal({
        isOpen,
        onClose,
        template,
}: EditTemplateTagsModalProps) {
        const router = useRouter();
        const [allTags, setAllTags] = useState<Tag[]>([]);
        const [selectedIDs, setSelectedIDs] = useState<number[]>(
                template.tags?.map((t) => t.id) ?? [],
        );
        const [isLoading, setIsLoading] = useState(false);

        useEffect(() => {
                if (!isOpen) return;
                setSelectedIDs(template.tags?.map((t) => t.id) ?? []);
                getTags().then((res) => {
                        if (res.success && res.data) setAllTags(res.data);
                });
        }, [isOpen]);

        const toggle = (id: number) => {
                setSelectedIDs((prev) => {
                        if (prev.includes(id)) return prev.filter((x) => x !== id);
                        if (prev.length >= MAX_TAGS) {
                                toast.error(`Poți selecta maximum ${MAX_TAGS} taguri.`);
                                return prev;
                        }
                        return [...prev, id];
                });
        };

        const handleConfirm = async () => {
                setIsLoading(true);
                const res = await setTemplateTags(template.id, selectedIDs);
                setIsLoading(false);
                if (!res.success) {
                        toast.error(res.message ?? "Eroare la salvarea tagurilor.");
                        return;
                }
                toast.success("Taguri actualizate!");
                router.refresh();
                onClose();
        };

        return (
                <Modal
                        isOpen={isOpen}
                        onClose={onClose}
                        onConfirm={handleConfirm}
                        title="Editează taguri"
                        closeOnConfirm={false}
                >
                        <p className="modal-subtext">
                                Selectează până la {MAX_TAGS} taguri pentru acest șablon.
                        </p>
                        <div className="tag-selector">
                                {allTags.map((tag) => {
                                        const selected = selectedIDs.includes(tag.id);
                                        const disabled =
                                                !selected && selectedIDs.length >= MAX_TAGS;
                                        return (
                                                <button
                                                        key={tag.id}
                                                        type="button"
                                                        className={`tag-option ${selected ? "tag-option--selected" : ""} ${disabled ? "tag-option--disabled" : ""}`}
                                                        style={
                                                                {
                                                                        "--tag-color":
                                                                                tag.color,
                                                                } as React.CSSProperties
                                                        }
                                                        onClick={() =>
                                                                !disabled && toggle(tag.id)
                                                        }
                                                >
                                                        {tag.name}
                                                        {selected && (
                                                                <span className="tag-option-check">
                                                                        ✓
                                                                </span>
                                                        )}
                                                </button>
                                        );
                                })}
                                {allTags.length === 0 && (
                                        <p className="modal-subtext">
                                                Nu există taguri disponibile.
                                        </p>
                                )}
                        </div>
                        <p className="tag-counter">
                                {selectedIDs.length} / {MAX_TAGS} selectate
                        </p>
                </Modal>
        );
}
