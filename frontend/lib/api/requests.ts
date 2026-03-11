"use server";
import { APIResponse, DocumentFile, DocumentRequest, PresignedURL, UserRole } from "@/types";
import { doclaneHTTPHelper } from "./core";

export async function getDocumentRequests(
        role: UserRole,
): Promise<APIResponse<DocumentRequest[]>> {
        return doclaneHTTPHelper(`/document-requests/${role.toLowerCase()}/my-requests`, {
                method: "GET",
        });
}

export async function getDocumentRequestById(
        requestId: string,
): Promise<APIResponse<DocumentRequest>> {
        return doclaneHTTPHelper(`/document-requests/${requestId}`, {
                method: "GET",
        });
}

export async function closeRequest(requestID: number): Promise<APIResponse> {
        return doclaneHTTPHelper(`/document-requests/${requestID}/archive`, {
                method: "POST",
                revalidate: "/dashboard/requests",
        });
}

export async function reopenRequest(requestID: number): Promise<APIResponse> {
        return doclaneHTTPHelper(`/document-requests/${requestID}/unarchive`, {
                method: "POST",
                revalidate: "/dashboard/requests",
        });
}

export async function createDocumentRequest(payload: {
        title: string;
        description?: string;
        client_id: number;
        is_recurring?: boolean;
        recurrence_cron?: string;
        is_scheduled?: boolean;
        scheduled_for?: string;
        due_date?: string;
        expected_documents: Array<{
                title: string;
                description: string;
                exampleFile?: File;
                exampleFileName?: string;
                ExampleMimeType?: string;
        }>;
}): Promise<APIResponse> {
        const hasExamples = payload.expected_documents.some((ed) => ed.exampleFile);

        if (hasExamples) {
                const formData = new FormData();
                formData.append("title", payload.title);
                if (payload.description) formData.append("description", payload.description);
                formData.append("client_id", payload.client_id.toString());
                if (payload.is_recurring) formData.append("is_recurring", "true");
                if (payload.recurrence_cron)
                        formData.append("recurrence_cron", payload.recurrence_cron);
                if (payload.is_scheduled) formData.append("is_scheduled", "true");
                if (payload.scheduled_for)
                        formData.append("scheduled_for", payload.scheduled_for);
                if (payload.due_date) formData.append("due_date", payload.due_date);

                payload.expected_documents.forEach((ed, i) => {
                        formData.append(`expected_documents[${i}][title]`, ed.title);
                        formData.append(
                                `expected_documents[${i}][description]`,
                                ed.description,
                        );
                        if (ed.exampleFile) {
                                formData.append(
                                        `expected_documents[${i}][example_file]`,
                                        ed.exampleFile,
                                );
                        }
                });

                return doclaneHTTPHelper("/document-requests", {
                        method: "POST",
                        formData,
                        revalidate: "/dashboard/requests",
                });
        }

        return doclaneHTTPHelper("/document-requests", {
                method: "POST",
                body: payload,
                revalidate: "/dashboard/requests",
        });
}

export async function presignDocumentURL(
        requestId: number,
        fileId: number,
): Promise<APIResponse<PresignedURL>> {
        return doclaneHTTPHelper(`/document-requests/${requestId}/files/${fileId}/presign`, {
                method: "GET",
        });
}

export async function presignExampleURL(
        expectedDocID: number,
): Promise<APIResponse<PresignedURL>> {
        return doclaneHTTPHelper(
                `/document-requests/expected-documents/${expectedDocID}/presign-example`,
                { method: "GET" },
        );
}

const MAX_FILE_SIZE = 20 * 1024 * 1024;

export async function uploadDocument(
        requestId: string,
        file: File,
        expectedDocumentId?: number,
): Promise<APIResponse> {
        const ALLOWED_EXTENSIONS = [".pdf", ".jpg", ".jpeg", ".png", ".doc", ".docx"];

        if (
                !ALLOWED_EXTENSIONS.includes(
                        file.name.substring(file.name.lastIndexOf(".")).toLowerCase(),
                )
        ) {
                throw new Error("File extension is not allowed.");
        }
        if (file.size > MAX_FILE_SIZE) {
                throw new Error("File exceeds the 20MB limit.");
        }

        const formData = new FormData();
        formData.append("file", file);
        if (expectedDocumentId !== undefined) {
                formData.append("expected_document_id", expectedDocumentId.toString());
        }

        return doclaneHTTPHelper(`/document-requests/${requestId}/files`, {
                method: "POST",
                formData,
                revalidate: `/dashboard/requests/${requestId}`,
        });
}

export async function updateExpectedDocumentStatus(
        expectedDocumentId: number,
        status: "approved" | "rejected" | "uploaded" | "pending",
        requestId: string,
        rejectionReason?: string,
): Promise<APIResponse> {
        return doclaneHTTPHelper(
                `/document-requests/expected-documents/${expectedDocumentId}/status`,
                {
                        method: "PATCH",
                        body: {
                                status,
                                ...(rejectionReason && { rejection_reason: rejectionReason }),
                        },
                        revalidate: `/dashboard/requests/${requestId}`,
                },
        );
}

export async function getFilesByRequestId(
        requestId: string,
): Promise<APIResponse<DocumentFile[]>> {
        return doclaneHTTPHelper(`/document-requests/${requestId}/files`, {
                method: "GET",
        });
}
