"use server";
import { APIResponse, User } from "@/types";
import { doclaneHTTPHelper } from "./core";

export async function getCurrentUser(): Promise<APIResponse<User>> {
        return doclaneHTTPHelper("/users/me", {
                method: "GET",
        });
}

export async function getUserById(userId: number): Promise<APIResponse<User>> {
        return doclaneHTTPHelper(`/users/${userId}`, {
                method: "GET",
        });
}

export async function getUsers(params?: {
        search?: string;
        limit?: number;
        offset?: number;
        orderBy?: string;
        order?: "asc" | "desc";
}): Promise<APIResponse<User[]>> {
        const query = new URLSearchParams();
        if (params?.search) query.append("search", params.search);
        if (params?.limit) query.append("limit", params.limit.toString());
        if (params?.offset) query.append("offset", params.offset.toString());
        if (params?.orderBy) query.append("order_by", params.orderBy);
        if (params?.order) query.append("order", params.order);

        const qs = query.toString();
        return doclaneHTTPHelper(`/users${qs ? `?${qs}` : ""}`, { method: "GET" });
}

export async function deactivateUser(userId: number): Promise<APIResponse> {
        return doclaneHTTPHelper(`/users/deactivate/${userId}`, {
                method: "POST",
                revalidate: "/dashboard/users",
        });
}

export async function notifyUser(userId: number): Promise<APIResponse> {
        return doclaneHTTPHelper(`/users/notify/${userId}`, {
                method: "POST",
        });
}

export async function getUsersByDepartment(
        departmentId: number,
): Promise<APIResponse<User[]>> {
        console.log(">>> getUsersByDepartment called with", departmentId);
        return doclaneHTTPHelper(`/users/by-department?department_id=${departmentId}`, {
                method: "GET",
        });
}

export async function updateUserDepartment(
        userId: number,
        departmentId: number,
): Promise<APIResponse> {
        return doclaneHTTPHelper(`/users/${userId}/department`, {
                method: "PATCH",
                body: { department_id: departmentId },
                revalidate: "/dashboard/departments",
        });
}

export async function updateUserProfile(payload: {
        phone?: string | null;
        street?: string | null;
        locality?: string | null;
}): Promise<APIResponse> {
        return doclaneHTTPHelper("/users/me/profile", {
                method: "PATCH",
                body: payload,
                revalidate: "/dashboard/settings",
        });
}
