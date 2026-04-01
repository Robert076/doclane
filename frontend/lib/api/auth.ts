"use server";
import { revalidatePath } from "next/cache";
import { cookies } from "next/headers";
import { APIResponse } from "@/types";
import { doclaneHTTPHelper } from "./core";

export async function login(email: string, password: string): Promise<APIResponse> {
        const data = await doclaneHTTPHelper("/auth/login", {
                method: "POST",
                body: { email, password },
        });

        if (!data.success) {
                return { success: false, message: data.message || "Login failed" };
        }

        const cookieStore = await cookies();
        cookieStore.set("auth_cookie", data.data as string, {
                httpOnly: true,
                secure: false,
                sameSite: "lax",
                path: "/",
                expires: new Date(Date.now() + 1000 * 60 * 60 * 24),
        });

        return { success: true, message: "Login successful." };
}

export async function logout(): Promise<APIResponse> {
        const cookieStore = await cookies();
        cookieStore.delete("auth_cookie");
        revalidatePath("/");
        return { success: true, message: "Logged out successfully" };
}

export async function register(
        email: string,
        password: string,
        invitationCode: string,
        firstName: string,
        lastName: string,
): Promise<APIResponse> {
        return doclaneHTTPHelper("/auth/register", {
                method: "POST",
                body: {
                        email,
                        password,
                        invitation_code: invitationCode,
                        first_name: firstName,
                        last_name: lastName,
                },
        });
}
