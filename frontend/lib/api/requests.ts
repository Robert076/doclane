"use server";
import {
        APIResponse,
        ALLOWED_EXTENSIONS,
        DocumentFile,
        Request,
        RequestComment,
} from "@/types";
import { doclaneHTTPHelper } from "./core";
import { cookies } from "next/headers";

export async function getAllRequests(search?: string): Promise<APIResponse<Request[]>> {
        const qs = search ? `?search=${encodeURIComponent(search)}` : "";
        return doclaneHTTPHelper(`/requests${qs}`, { method: "GET" });
}

export async function getRequestsByAssignee(
        assigneeId: number,
): Promise<APIResponse<Request[]>> {
        return doclaneHTTPHelper(`/requests/assignee/${assigneeId}`, {
                method: "GET",
        });
}

export async function getRequestsByDepartment(
        departmentId: number,
): Promise<APIResponse<Request[]>> {
        return doclaneHTTPHelper(`/requests/department/${departmentId}`, {
                method: "GET",
        });
}

export async function getRequestById(requestId: number): Promise<APIResponse<Request>> {
        return doclaneHTTPHelper(`/requests/${requestId}`, {
                method: "GET",
        });
}

export async function createRequest(payload: {
        template_id: number;
        is_scheduled?: boolean;
        scheduled_for?: string;
        due_date?: string;
}): Promise<APIResponse<number>> {
        return doclaneHTTPHelper("/requests", {
                method: "POST",
                body: payload,
                revalidate: "/dashboard/requests",
        });
}

export async function forwardRequestToDepartment(
        requestId: number,
        departmentId: number,
): Promise<APIResponse> {
        return doclaneHTTPHelper(`/requests/forward/${requestId}`, {
                method: "POST",
                body: { department_id: departmentId },
                revalidate: `/dashboard/requests/${requestId}`,
        });
}

export async function getArchivedRequests(): Promise<APIResponse<Request[]>> {
        return doclaneHTTPHelper("/requests/archived", { method: "GET" });
}

export async function getCancelledRequests(): Promise<APIResponse<Request[]>> {
        return doclaneHTTPHelper("/requests/cancelled", { method: "GET" });
}

export async function claimRequest(requestId: number): Promise<APIResponse> {
        return doclaneHTTPHelper(`/requests/${requestId}/claim`, {
                method: "POST",
                revalidate: "/dashboard/requests",
        });
}

export async function unclaimRequest(requestId: number): Promise<APIResponse> {
        return doclaneHTTPHelper(`/requests/${requestId}/unclaim`, {
                method: "POST",
                revalidate: "/dashboard/requests",
        });
}

export async function closeRequest(requestID: number): Promise<APIResponse> {
        return doclaneHTTPHelper(`/requests/${requestID}/archive`, {
                method: "POST",
                revalidate: "/dashboard/requests",
        });
}

export async function reopenRequest(requestID: number): Promise<APIResponse> {
        return doclaneHTTPHelper(`/requests/${requestID}/unarchive`, {
                method: "POST",
                revalidate: "/dashboard/requests",
        });
}

export async function cancelRequest(requestId: number): Promise<APIResponse> {
        return doclaneHTTPHelper(`/requests/${requestId}/cancel`, {
                method: "POST",
                revalidate: "/dashboard/requests",
        });
}

export async function addComment(requestId: number, comment: string): Promise<APIResponse> {
        return doclaneHTTPHelper(`/requests/${requestId}/comments`, {
                method: "POST",
                body: { comment },
                revalidate: `/dashboard/requests/${requestId}`,
        });
}

export async function getCommentsByRequest(
        requestId: number,
): Promise<APIResponse<RequestComment[]>> {
        return doclaneHTTPHelper(`/requests/${requestId}/comments`, {
                method: "GET",
        });
}

export async function presignDocumentURL(
        requestId: number,
        fileId: number,
): Promise<APIResponse<string>> {
        return doclaneHTTPHelper(`/requests/${requestId}/files/${fileId}/presign`, {
                method: "GET",
        });
}

export async function extractFileText(
        requestId: number,
        fileId: number,
): Promise<APIResponse<{ text: string }>> {
        return doclaneHTTPHelper(`/requests/${requestId}/files/${fileId}/extract`, {
                method: "GET",
        });
}

export async function interpretFileText(
        requestId: number,
        fileId: number,
        documentTitle: string,
): Promise<APIResponse<{ interpretation: string }>> {
        return doclaneHTTPHelper(
                `/requests/${requestId}/files/${fileId}/interpret?title=${encodeURIComponent(documentTitle)}`,
                { method: "GET" },
        );
}

// in requests.ts
export async function speakFileText(
        requestId: number,
        fileId: number,
): Promise<APIResponse<{ audio: string }>> {
        const cookieStore = await cookies();
        const authCookie = cookieStore.get("auth_cookie");
        const fetchUrl = `${process.env.BACKEND_URL}/requests/${requestId}/files/${fileId}/speak`;

        try {
                const response = await fetch(fetchUrl, {
                        method: "GET",
                        headers: {
                                Authorization: `Bearer ${authCookie?.value}`,
                        },
                });

                if (!response.ok) {
                        return {
                                success: false,
                                message: "Eroare la generarea audio.",
                                error: "Failed",
                        };
                }

                const buffer = await response.arrayBuffer();
                const base64 = Buffer.from(buffer).toString("base64");

                return { success: true, message: "OK", data: { audio: base64 } };
        } catch (error) {
                logger.error(`Error during speak: ${error}`);
                return {
                        success: false,
                        message: "Something went wrong",
                        error: "Something went wrong",
                };
        }
}

export async function presignExampleURL(expectedDocID: number): Promise<APIResponse<string>> {
        const res = await doclaneHTTPHelper<{ url: string }>(
                `/requests/expected-documents/${expectedDocID}/presign-example`,
                { method: "GET" },
        );
        return { ...res, data: res.data?.url };
}

const MAX_FILE_SIZE = 20 * 1024 * 1024;

export async function uploadDocument(
        requestId: number,
        file: File,
        expectedDocumentId?: number,
): Promise<APIResponse> {
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

        return doclaneHTTPHelper(`/requests/${requestId}/files`, {
                method: "POST",
                formData,
                revalidate: `/dashboard/requests/${requestId}`,
        });
}

export async function updateExpectedDocumentStatus(
        expectedDocumentId: number,
        status: "accepted" | "rejected" | "uploaded" | "pending",
        requestId: number,
        rejectionReason?: string,
): Promise<APIResponse> {
        return doclaneHTTPHelper(`/requests/expected-documents/${expectedDocumentId}/status`, {
                method: "PATCH",
                body: {
                        status,
                        ...(rejectionReason && { rejection_reason: rejectionReason }),
                },
                revalidate: `/dashboard/requests/${requestId}`,
        });
}

export async function getFilesByRequestId(
        requestId: number,
): Promise<APIResponse<DocumentFile[]>> {
        return doclaneHTTPHelper(`/requests/${requestId}/files`, {
                method: "GET",
        });
}

export async function patchRequest(
        requestId: number,
        payload: { title: string },
): Promise<APIResponse> {
        return doclaneHTTPHelper(`/requests/${requestId}`, {
                method: "PATCH",
                body: payload,
                revalidate: `/dashboard/requests/${requestId}`,
        });
}
