import { updateExpectedDocumentStatus } from "@/lib/api/requests";
import { useState } from "react";
import toast from "react-hot-toast";
import { ExpectedDocumentStatus } from "@/types";

export function useDocumentStatus(
        expectedDocumentId: number,
        requestId: number,
        hasFiles: boolean,
) {
        const [isLoading, setIsLoading] = useState(false);

        const updateStatus = async (status: ExpectedDocumentStatus, reason?: string) => {
                setIsLoading(true);
                const res = await updateExpectedDocumentStatus(
                        expectedDocumentId,
                        status,
                        requestId,
                        reason,
                );
                setIsLoading(false);
                if (!res.success) {
                        toast.error(
                                res.message ?? "Nu s-a putut actualiza statusul documentului.",
                        );
                }
        };

        const approve = () => updateStatus("accepted");
        const reject = (reason: string) => updateStatus("rejected", reason);
        const reset = () => updateStatus(hasFiles ? "uploaded" : "pending");

        return { approve, reject, reset, isLoading };
}
