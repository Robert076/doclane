import { useState } from "react";
import { updateExpectedDocumentStatus } from "@/lib/api/api";

type DocumentStatus = "approved" | "rejected" | "uploaded" | "pending";

export function useDocumentStatus(
        expectedDocumentId: string,
        requestId: string,
        hasFiles: boolean,
) {
        const [isLoading, setIsLoading] = useState(false);

        const updateStatus = async (status: DocumentStatus, reason?: string) => {
                setIsLoading(true);
                await updateExpectedDocumentStatus(
                        +expectedDocumentId,
                        status,
                        requestId,
                        reason,
                );
                setIsLoading(false);
        };

        const approve = () => updateStatus("approved");
        const reject = (reason: string) => updateStatus("rejected", reason);
        const reset = () => updateStatus(hasFiles ? "uploaded" : "pending");

        return { approve, reject, reset, isLoading };
}
