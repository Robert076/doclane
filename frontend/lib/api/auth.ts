"use server";
import { revalidatePath } from "next/cache";
import { cookies } from "next/headers";
import { APIResponse } from "@/types";
import { doclaneHTTPHelper } from "./core";

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

export async function logout(): Promise<APIResponse> {
        const cookieStore = await cookies();
        cookieStore.delete("auth_cookie");
        revalidatePath("/");

        return {
                success: true,
                message: "Logged out successfully",
        };
}

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
                        invitation_code: invitationCode,
                        first_name: firstName,
                        last_name: lastName,
                },
        });
}
