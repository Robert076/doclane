"use server";
import { APIResponse } from "@/types";
import { doclaneHTTPHelper } from "./core";
import { Stats } from "@/types/stats";

export async function getStats(): Promise<APIResponse<Stats>> {
        return doclaneHTTPHelper("/stats", { method: "GET" });
}
