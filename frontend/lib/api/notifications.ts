// Add to lib/api/notifications.ts

"use server";
import { APIResponse, AuditEvent } from "@/types";
import { doclaneHTTPHelper } from "./core";

export async function getNotifications(limit?: number): Promise<APIResponse<AuditEvent[]>> {
  const qs = limit ? `?limit=${limit}` : "";
  return doclaneHTTPHelper(`/notifications${qs}`, {
    method: "GET",
  });
}

export async function markNotificationsSeen(): Promise<APIResponse> {
    return doclaneHTTPHelper("/notifications/seen", {
        method: "POST",
    });
}