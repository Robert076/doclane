"use server";
import { revalidatePath } from "next/cache";
import { cookies } from "next/headers";
import { logger } from "../logger";

interface APIResponse {
        success: boolean;
        message: string;
        error?: string;
        data?: any;
}

const BACKEND_URL = process.env.BACKEND_URL!;

export async function doclaneHTTPHelper(
        fetchUrl: string,
        method: string,
): Promise<APIResponse> {
        try {
                const cookieStore = await cookies();
                const authCookie = cookieStore.get("auth_cookie");
                const response = await fetch(fetchUrl, {
                        method: method,
                        credentials: "include",
                        headers: {
                                Authorization: `Bearer ${authCookie?.value}`,
                        },
                });

                const resultData = await response.json();

                if (!response.ok) {
                        const errorData = await response.json();
                        logger.error(`Error during HTTP call: ${errorData}`);
                        return {
                                success: false,
                                error: errorData,
                                message: resultData.message,
                        };
                }

                logger.info(`${method} ${fetchUrl} call successful`);

                return {
                        success: true,
                        message: resultData.message,
                        data: resultData.data,
                };
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
        const fetchUrl = `${BACKEND_URL}/users/deactivate/${userId}`;
        const method = "POST";

        const result = await doclaneHTTPHelper(fetchUrl, method);
        if (result.success === false) {
                return {
                        success: false,
                        message:
                                "Error occured when trying to deactivate user: " +
                                result.message,
                        error: result.error,
                };
        }

        revalidatePath("/dashboard/clients");
        return {
                success: true,
                message: "User deactivated successfully.",
        };
}

export async function presignDocumentURL(
        requestId: number,
        fileId: number,
): Promise<APIResponse> {
        const fetchUrl = `${BACKEND_URL}/document-requests/${requestId}/files/${fileId}/presign`;
        const method = "GET";

        const result = await doclaneHTTPHelper(fetchUrl, method);
        if (result.success === false) {
                return {
                        success: false,
                        message:
                                "Error occured when trying to presign file: " + result.message,
                        error: result.error,
                };
        }

        return {
                success: true,
                message: result.message,
                data: result.data,
        };
}
