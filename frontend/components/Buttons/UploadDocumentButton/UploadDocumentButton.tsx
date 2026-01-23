"use client";

import { useRef, useState } from "react";
import ButtonPrimary from "@/components/Buttons/ButtonPrimary/ButtonPrimary";
import { useRouter } from "next/navigation";
import { uploadDocument, ALLOWED_EXTENSIONS } from "@/lib/uploadDocument";
import { toast } from "react-hot-toast";

interface Props {
  requestId: string;
}

export default function UploadDocumentButton({ requestId }: Props) {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [isUploading, setIsUploading] = useState(false);
  const router = useRouter();

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    const resetInput = () => {
      if (fileInputRef.current) fileInputRef.current.value = "";
    };

    setIsUploading(true);

    toast.promise(
      uploadDocument(requestId, file)
        .then((res) => {
          router.refresh();
          setIsUploading(false);
          resetInput();
          return res;
        })
        .catch((err) => {
          setIsUploading(false);
          resetInput();
          throw err;
        }),
      {
        loading: "Uploading document...",
        success: "Document uploaded successfully!",
        error: (err) => `Error: ${err.message}`,
      }
    );
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
        text={isUploading ? "Se încarcă..." : "Upload Document"}
        variant="primary"
        fullWidth
        disabled={isUploading}
        onClick={() => fileInputRef.current?.click()}
      />
    </>
  );
}
