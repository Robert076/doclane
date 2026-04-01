import { updateExpectedDocumentStatus } from "@/lib/api/requests";
import { useState } from "react";
import { ExpectedDocumentStatus } from "@/types";

export function useDocumentStatus(
        expectedDocumentId: number,
        requestId: number,
        hasFiles: boolean,
) {
        const [isLoading, setIsLoading] = useState(false);

        const updateStatus = async (status: ExpectedDocumentStatus, reason?: string) => {
                setIsLoading(true);
                await updateExpectedDocumentStatus(
                        expectedDocumentId,
                        status,
                        requestId,
                        reason,
                );
                setIsLoading(false);
        };

        const approve = () => updateStatus("accepted");
        const reject = (reason: string) => updateStatus("rejected", reason);
        const reset = () => updateStatus(hasFiles ? "uploaded" : "pending");

        return { approve, reject, reset, isLoading };
}
