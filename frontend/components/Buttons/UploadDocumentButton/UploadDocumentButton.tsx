"use client";

import { useRef, useState } from "react";
import ButtonPrimary from "@/components/Buttons/ButtonPrimary/ButtonPrimary";
import { useRouter } from "next/navigation";

interface Props {
  requestId: string;
}

// Extensiile permise (la fel ca in Go)
const ALLOWED_EXTENSIONS = [".pdf", ".jpg", ".jpeg", ".png", ".doc", ".docx"];
const MAX_FILE_SIZE = 20 * 1024 * 1024; // 20MB

export default function UploadDocumentButton({ requestId }: Props) {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [isUploading, setIsUploading] = useState(false);
  const router = useRouter();

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // --- Validări Client-Side (Oglindesc backend-ul Go) ---
    const fileExt = file.name.substring(file.name.lastIndexOf(".")).toLowerCase();
    if (!ALLOWED_EXTENSIONS.includes(fileExt)) {
      alert(`Extensia ${fileExt} nu este permisă.`);
      return;
    }

    if (file.size > MAX_FILE_SIZE) {
      alert("Fișierul depășește limita de 20MB.");
      return;
    }

    // --- Pregătirea FormData ---
    const formData = new FormData();
    formData.append("file", file);

    setIsUploading(true);

    try {
      // Ajustează URL-ul conform rutei tale din backend
      const response = await fetch(
        `http://localhost:8080/api/document-requests/${requestId}/files`,
        {
          method: "POST",
          body: formData,
          credentials: "include",
          // Atenție: NU seta Content-Type manual când trimiți FormData,
          // browser-ul o va face automat cu tot cu boundary string.
        }
      );

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.msg || "Upload failed");
      }

      alert("Document încărcat cu succes!");
      router.refresh(); // Reîmprospătează datele în pagină
    } catch (error: any) {
      alert(`Eroare: ${error.message}`);
    } finally {
      setIsUploading(false);
      if (fileInputRef.current) fileInputRef.current.value = "";
    }
  };

  return (
    <>
      <input
        type="file"
        ref={fileInputRef}
        onChange={handleFileChange}
        style={{ display: "none" }}
        accept={ALLOWED_EXTENSIONS.join(",")}
      />
      <ButtonPrimary
        text={isUploading ? "Uploading..." : "Upload Document"}
        variant="primary"
        fullWidth
        disabled={isUploading}
        onClick={() => fileInputRef.current?.click()}
      />
    </>
  );
}
