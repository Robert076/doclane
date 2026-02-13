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

export async function deactivateUser(userId: number): Promise<APIResponse> {
        try {
                const cookieStore = await cookies();
                const authCookie = cookieStore.get("auth_cookie");
                const res = await fetch(`${BACKEND_URL}/users/deactivate/${userId}`, {
                        method: "POST",
                        credentials: "include",
                        headers: {
                                Authorization: `Bearer ${authCookie?.value}`,
                        },
                });

                if (!res.ok) {
                        const errorData = await res.json();
                        logger.error(`Error deactivating client: ${errorData}`);
                        return {
                                success: false,
                                error: errorData,
                                message: "An error occured when attempting to deactivate user.",
                        };
                }

                revalidatePath("/dashboard/clients");
                logger.info(`Successfully deactivated client with id ${userId}`);

                return {
                        success: true,
                        message: "User deactivated successfully",
                };
        } catch (error) {
                logger.error(`Error deactivating client: ${error}`);
                return {
                        success: false,
                        message: "Something went wrong",
                        error: "Something went wrong",
                };
        }
}
