"use server";
import { APIResponse, User } from "@/types";
import { doclaneHTTPHelper } from "./core";
import { getRequestById } from "./requests";

export async function getCurrentUser(): Promise<APIResponse<User>> {
        return doclaneHTTPHelper("/users/me", {
                method: "GET",
        });
}

export async function getUserById(userId: string): Promise<APIResponse> {
        return doclaneHTTPHelper(`/users/${userId}`, {
                method: "GET",
        });
}

export async function getClientsByProfessional(): Promise<APIResponse<User[]>> {
        return doclaneHTTPHelper(`/users/my-clients`, { method: "GET" });
}

export async function deactivateUser(userId: number): Promise<APIResponse> {
        return doclaneHTTPHelper(`/users/deactivate/${userId}`, {
                method: "POST",
                revalidate: "/dashboard/clients",
        });
}

export async function sendEmail(requestId: number): Promise<APIResponse> {
        const responseRequest = await getRequestById(requestId.toString());
        if (responseRequest.success === false || !responseRequest.data) {
                return responseRequest;
        }

        const request = responseRequest.data;

        const response = await doclaneHTTPHelper(`/users/notify/${request.client_id}`, {
                method: "POST",
        });

        return response;
}
