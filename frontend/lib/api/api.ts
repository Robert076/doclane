"use server";
import { revalidatePath } from "next/cache";
import { cookies } from "next/headers";
import { logger } from "../logger";
import { DocumentRequest, UserRole } from "@/types";

interface APIResponse {
        success: boolean;
        message: string;
        error?: string;
        data?: any;
}

interface HTTPOptions {
        method?: "GET" | "POST" | "PUT" | "DELETE" | "PATCH";
        body?: any;
        formData?: FormData;
        revalidate?: string;
}

const BACKEND_URL = process.env.BACKEND_URL!;

export async function doclaneHTTPHelper(
        endpoint: string,
        options: HTTPOptions,
): Promise<APIResponse> {
        const { method = "GET", body, formData, revalidate } = options;
        const fetchUrl = `${BACKEND_URL}${endpoint}`;

        try {
                const cookieStore = await cookies();
                const authCookie = cookieStore.get("auth_cookie");
                const response = await fetch(fetchUrl, {
                        method,
                        credentials: "include",
                        headers: {
                                Authorization: `Bearer ${authCookie?.value}`,
                                ...(body && { "Content-Type": "application/json" }),
                        },
                        ...(formData
                                ? { body: formData }
                                : body
                                  ? { body: JSON.stringify(body) }
                                  : {}),
                });

                const resultData = await response.json();

                if (!response.ok) {
                        logger.error(
                                `Error during ${method} ${fetchUrl}: ${resultData.message || resultData.error}`,
                        );
                        return {
                                success: false,
                                error:
                                        resultData.error ||
                                        resultData.message ||
                                        "Request failed",
                                message:
                                        resultData.message ||
                                        resultData.error ||
                                        "Request failed",
                        };
                }

                logger.info(`${method} ${fetchUrl} call successful`);
                if (revalidate) revalidatePath(revalidate);

                return { success: true, message: resultData.message, data: resultData.data };
        } catch (error) {
                logger.error(`Error during HTTP call: ${error}`);
                return {
                        success: false,
                        message: "Something went wrong",
                        error: "Something went wrong",
                };
        }
}

export async function deactivateUser(userId: number): Promise<APIResponse> {
        return doclaneHTTPHelper(`/users/deactivate/${userId}`, {
                method: "POST",
                revalidate: "/dashboard/clients",
        });
}

export async function presignDocumentURL(
        requestId: number,
        fileId: number,
): Promise<APIResponse> {
        return doclaneHTTPHelper(`/document-requests/${requestId}/files/${fileId}/presign`, {
                method: "GET",
        });
}

export async function logout(): Promise<APIResponse> {
        // logout does not use the backend since it's redundant
        const cookieStore = await cookies();
        cookieStore.delete("auth_cookie");
        revalidatePath("/");

        return {
                success: true,
                message: "Logged out successfully",
        };
}

export async function getDocumentRequests(role: UserRole): Promise<APIResponse> {
        return doclaneHTTPHelper(`/document-requests/${role.toLowerCase()}/my-requests`, {
                method: "GET",
        });
}

export async function getCurrentUser(): Promise<APIResponse> {
        return doclaneHTTPHelper("/users/me", {
                method: "GET",
        });
}

export async function getDocumentRequestById(requestId: string): Promise<APIResponse> {
        return doclaneHTTPHelper(`/document-requests/${requestId}`, {
                method: "GET",
        });
}

export async function getUserById(userId: string): Promise<APIResponse> {
        return doclaneHTTPHelper(`/users/${userId}`, {
                method: "GET",
        });
}

export async function sendEmail(requestId: number): Promise<APIResponse> {
        const responseRequest = await getDocumentRequestById(requestId.toString());
        if (responseRequest.success === false) {
                return responseRequest;
        }

        const request: DocumentRequest = responseRequest.data;
        const response = await doclaneHTTPHelper(`/users/notify/${request.client_id}`, {
                method: "POST",
        });

        return response;
}

export async function closeRequest(requestID: number): Promise<APIResponse> {
        return doclaneHTTPHelper(`/document-requests/${requestID}/deactivate`, {
                method: "POST",
        });
}

export async function createDocumentRequest(payload: object): Promise<APIResponse> {
        return doclaneHTTPHelper("/document-requests", {
                method: "POST",
                body: payload,
                revalidate: "/dashboard/requests",
        });
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
