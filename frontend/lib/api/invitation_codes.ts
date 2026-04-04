"use server";
import { APIResponse, InvitationCode } from "@/types";
import { doclaneHTTPHelper } from "./core";

export async function generateInvitationCode(
        departmentId: number,
        expiresInDays?: number,
): Promise<APIResponse<{ code: string }>> {
        return doclaneHTTPHelper("/invitations/generate", {
                method: "POST",
                body: {
                        department_id: departmentId,
                        expires_in_days: expiresInDays ?? 7,
                },
        });
}

export async function getInvitationCodesByDepartment(
        departmentId: number,
): Promise<APIResponse<InvitationCode[]>> {
        return doclaneHTTPHelper(`/invitations/by-department?department_id=${departmentId}`, {
                method: "GET",
        });
}

export async function deleteInvitationCode(id: number): Promise<APIResponse> {
        return doclaneHTTPHelper(`/invitations/${id}`, {
                method: "DELETE",
                revalidate: "/dashboard/departments",
        });
}
