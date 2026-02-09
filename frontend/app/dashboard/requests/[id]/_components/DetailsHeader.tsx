"use client";

import { DocumentRequest } from "@/types";
import Link from "next/link";
import "./DetailsHeader.css";
import { MdEdit, MdCheck, MdClose } from "react-icons/md";
import { useState } from "react";
import toast from "react-hot-toast";

export default function DetailsHeader({ data }: { data: DocumentRequest }) {
  const [isEditing, setIsEditing] = useState(false);
  const [title, setTitle] = useState(data.title);

  const handleSave = async () => {
    if (!title.trim()) {
      toast.error("Title cannot be empty");
      return;
    }

    const updatePromise = fetch(`/api/backend/document-requests/${data.id}`, {
      method: "PATCH",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ title }),
    }).then(async (res) => {
      if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData.error || "Failed to update title");
      }
      return res.json();
    });

    toast.promise(updatePromise, {
      loading: "Updating title...",
      success: "Title updated successfully!",
      error: (err) => `Failed: ${err.message}`,
    });

    updatePromise
      .then(() => {
        setIsEditing(false);
      })
      .catch(() => {
        setTitle(data.title);
      });
  };

  const handleCancel = () => {
    setTitle(data.title);
    setIsEditing(false);
  };

  return (
    <header className="details-header">
      <Link href="/dashboard" className="back-link">
        ← Înapoi la Cereri
      </Link>
      <div className="header-main">
        {isEditing ? (
          <div className="title-edit-container">
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="title-input"
              autoFocus
              onKeyDown={(e) => {
                if (e.key === "Enter") handleSave();
                if (e.key === "Escape") handleCancel();
              }}
            />
            <button onClick={handleSave} className="icon-button save-button">
              <MdCheck />
            </button>
            <button onClick={handleCancel} className="icon-button cancel-button">
              <MdClose />
            </button>
          </div>
        ) : (
          <div className="title-display-container">
            <h1>{title}</h1>
            <button onClick={() => setIsEditing(true)} className="icon-button edit-button">
              <MdEdit />
            </button>
          </div>
        )}
      </div>
    </header>
  );
}
