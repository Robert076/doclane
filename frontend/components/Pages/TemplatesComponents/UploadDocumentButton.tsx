"use client";
import { useRef, useState } from "react";
import ButtonPrimary from "@/components/ButtonComponents/ButtonPrimary/ButtonPrimary";
import { useRouter } from "next/navigation";

import { toast } from "react-hot-toast";
import { uploadDocument } from "@/lib/api/api";
import { UI_TEXT } from "@/locales/ro";

const ALLOWED_EXTENSIONS = [".pdf", ".jpg", ".jpeg", ".png", ".doc", ".docx"];
interface Props {
        requestId: string;
        expectedDocumentId?: number;
}

export default function UploadDocumentButton({ requestId, expectedDocumentId }: Props) {
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
                        uploadDocument(requestId, file, expectedDocumentId)
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
                        },
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
                                text={
                                        isUploading
                                                ? UI_TEXT.buttons.uploadDocument.inProgress
                                                : UI_TEXT.buttons.uploadDocument.normal
                                }
                                variant="primary"
                                fullWidth
                                disabled={isUploading}
                                onClick={() => fileInputRef.current?.click()}
                        />
                </>
        );
}
