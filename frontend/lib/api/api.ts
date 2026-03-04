"use server";
import { revalidatePath } from "next/cache";
import { cookies } from "next/headers";
import { logger } from "../logger";
import { APIResponse, DocumentRequest, UserRole } from "@/types";

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
        return doclaneHTTPHelper(`/document-requests/${requestID}/archive`, {
                method: "POST",
                revalidate: "/dashboard",
        });
}

export async function archiveTemplate(id: number): Promise<APIResponse> {
        return doclaneHTTPHelper(`/templates/${id}/archive`, {
                method: "POST",
                revalidate: "/dashboard/templates",
        });
}

export async function unarchiveTemplate(id: number): Promise<APIResponse> {
        return doclaneHTTPHelper(`/templates/${id}/unarchive`, {
                method: "POST",
                revalidate: "/dashboard/templates",
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

export async function login(email: string, password: string) {
        const res = await fetch(`${process.env.NEXT_PUBLIC_APP_URL}/api/login`, {
                method: "POST",
                headers: {
                        "Content-Type": "application/json",
                },
                body: JSON.stringify({ email, password }),
        });

        return res.json();
}

export async function presignExampleURL(expectedDocID: number): Promise<APIResponse> {
        return doclaneHTTPHelper(
                `/document-requests/expected-documents/${expectedDocID}/presign-example`,
                { method: "GET" },
        );
}

export async function presignTemplateExample(
        templateID: number,
        docID: number,
): Promise<APIResponse> {
        return doclaneHTTPHelper(
                `/templates/${templateID}/expected-documents/${docID}/presign-example`,
                { method: "GET" },
        );
}

export async function getDocumentRequestTemplateByID(id: number): Promise<APIResponse> {
        return doclaneHTTPHelper(`/templates/${id}`, { method: "GET" });
}

export async function getClientsByProfessional(): Promise<APIResponse> {
        return doclaneHTTPHelper(`/users/my-clients`, { method: "GET" });
}

// const loginPromise = fetch("/api/backend/auth/login", {
//                         method: "POST",
//                         credentials: "include",
//                         headers: {
//                                 "Content-Type": "application/json",
//                         },
//                         body: JSON.stringify({ email, password }),
//                 }).then(async (res) => {
//                         if (!res.ok) {
//                                 const errorData = await res.json();
//                                 throw new Error(errorData.error || "Login failed");
//                         }
//                         return res.json();
//                 });

//                 toast.promise(loginPromise, {
//                         loading: "Logging in...",
//                         success: "Login successful!",
//                         error: (err) => `Login failed: ${err.message}`,
//                 });

//                 loginPromise.then((_) => {
//                         router.push("/dashboard");
//                 });

export async function signUpClient(
        email: string,
        password: string,
        invitationCode: string,
        firstName: string,
        lastName: string,
) {
        return doclaneHTTPHelper("/auth/register/client", {
                method: "POST",
                body: {
                        email,
                        password,
                        invitationCode,
                        firstName,
                        lastName,
                },
        });
}

// fetch("/api/backend/auth/register/client", {
//                         method: "POST",
//                         credentials: "include",
//                         headers: {
//                                 "Content-Type": "application/json",
//                         },
//                         body: JSON.stringify({
//                                 email,
//                                 password,
//                                 invitation_code: invitationCode,
//                                 first_name: firstName,
//                                 last_name: lastName,
//                         }),
//                 }).then(async (res) => {
//                         if (!res.ok) {
//                                 const errorData = await res.json();
//                                 throw new Error(errorData.error || "Sign up failed");
//                         }
//                         return res.json();
//                 });

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

export async function getTemplates(): Promise<APIResponse> {
        return doclaneHTTPHelper("/templates", {
                method: "GET",
        });
}

export async function getTemplateByID(templateID: number): Promise<APIResponse> {
        return doclaneHTTPHelper(`/templates/${templateID}`, {
                method: "GET",
        });
}

export async function createTemplate(payload: object): Promise<APIResponse> {
        return doclaneHTTPHelper("/templates", {
                method: "POST",
                body: payload,
                revalidate: "/dashboard/templates",
        });
}

export async function addExpectedDocumentTemplate(
        templateID: number,
        title: string,
        description: string,
        exampleFile?: File,
): Promise<APIResponse> {
        const formData = new FormData();
        formData.append("title", title);
        formData.append("description", description);
        if (exampleFile) {
                formData.append("example_file", exampleFile);
        }

        return doclaneHTTPHelper(`/templates/${templateID}/expected-documents`, {
                method: "POST",
                formData,
                revalidate: `/dashboard/templates/${templateID}`,
        });
}

export async function getExpectedDocumentTemplatesByTemplate(
        templateID: number,
): Promise<APIResponse> {
        return doclaneHTTPHelper(`/templates/${templateID}/expected-documents`, {
                method: "GET",
        });
}

export async function deleteExpectedDocumentTemplate(
        templateID: number,
        expectedDocTemplateID: number,
): Promise<APIResponse> {
        return doclaneHTTPHelper(
                `/templates/${templateID}/expected-documents/${expectedDocTemplateID}`,
                {
                        method: "DELETE",
                        revalidate: `/dashboard/templates/${templateID}`,
                },
        );
}

export async function instantiateTemplate(
        templateID: number,
        payload: {
                client_id: number;
                is_scheduled: boolean;
                scheduled_for?: string;
                due_date?: string;
        },
): Promise<APIResponse> {
        return doclaneHTTPHelper(`/templates/${templateID}/instantiate`, {
                method: "POST",
                body: payload,
                revalidate: "/dashboard/requests",
        });
}
