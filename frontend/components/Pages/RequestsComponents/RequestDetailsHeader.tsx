"use client";
import { Request } from "@/types";
import Link from "next/link";
import { MdEdit, MdCheck, MdClose } from "react-icons/md";
import { useState } from "react";
import toast from "react-hot-toast";
import { useUser } from "@/context/UserContext";
import { patchRequest } from "@/lib/api/requests";
import "./RequestDetailsHeader.css";

export default function RequestDetailsHeader({ data }: { data: Request }) {
        const user = useUser();
        const [isEditing, setIsEditing] = useState(false);
        const [title, setTitle] = useState(data.title);

        const canEdit = user.role === "admin" || user.department_id !== null;

        const handleSave = async () => {
                if (!title.trim()) {
                        toast.error("Titlul nu poate fi gol.");
                        return;
                }

                const response = await patchRequest(data.id, { title });
                if (response.success) {
                        toast.success("Titlul a fost actualizat.");
                        setIsEditing(false);
                } else {
                        toast.error(response.message);
                        setTitle(data.title);
                }
        };

        const handleCancel = () => {
                setTitle(data.title);
                setIsEditing(false);
        };

        return (
                <header className="details-header">
                        <Link href="/dashboard/requests" className="back-link">
                                ← Înapoi la Cereri
                        </Link>
                        <div className="header-main">
                                {isEditing ? (
                                        <div className="title-edit-container">
                                                <input
                                                        type="text"
                                                        value={title}
                                                        onChange={(e) =>
                                                                setTitle(e.target.value)
                                                        }
                                                        className="title-input"
                                                        autoFocus
                                                        onKeyDown={(e) => {
                                                                if (e.key === "Enter")
                                                                        handleSave();
                                                                if (e.key === "Escape")
                                                                        handleCancel();
                                                        }}
                                                />
                                                <button
                                                        onClick={handleSave}
                                                        className="icon-button save-button"
                                                >
                                                        <MdCheck />
                                                </button>
                                                <button
                                                        onClick={handleCancel}
                                                        className="icon-button cancel-button"
                                                >
                                                        <MdClose />
                                                </button>
                                        </div>
                                ) : (
                                        <div className="title-display-container">
                                                <h1>{title}</h1>
                                                {canEdit && (
                                                        <button
                                                                onClick={() =>
                                                                        setIsEditing(true)
                                                                }
                                                                className="icon-button edit-button"
                                                        >
                                                                <MdEdit />
                                                        </button>
                                                )}
                                        </div>
                                )}
                        </div>
                </header>
        );
}
