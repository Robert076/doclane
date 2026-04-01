"use server";
import { APIResponse, Department } from "@/types";
import { doclaneHTTPHelper } from "./core";

export async function getDepartments(): Promise<APIResponse<Department[]>> {
        return doclaneHTTPHelper("/departments", { method: "GET" });
}

export async function createDepartment(name: string): Promise<APIResponse<number>> {
        return doclaneHTTPHelper("/departments", {
                method: "POST",
                body: { name },
                revalidate: "/dashboard/settings",
        });
}
