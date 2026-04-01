"use server";

import { revalidatePath } from "next/cache";
import { cookies } from "next/headers";
import { logger } from "../logger";
import { APIResponse } from "@/types";

interface HTTPOptions {
        method?: "GET" | "POST" | "PUT" | "DELETE" | "PATCH";
        body?: any;
        formData?: FormData;
        revalidate?: string;
}

const BACKEND_URL = process.env.BACKEND_URL!;

export async function doclaneHTTPHelper<T = unknown>(
        endpoint: string,
        options: HTTPOptions,
): Promise<APIResponse<T>> {
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
