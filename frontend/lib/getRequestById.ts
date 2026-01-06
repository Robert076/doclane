import { ApiResponse, DocumentRequest } from "@/types";
import { cookies } from "next/headers";

export type RequestFetchResult = {
  data: DocumentRequest | null;
  error?: "FORBIDDEN" | "NOT_FOUND" | "SERVER_ERROR";
};

export default async function getRequestById(id: string): Promise<RequestFetchResult> {
  const cookieStore = await cookies();
  const token = cookieStore.get("auth_cookie")?.value;

  try {
    const res = await fetch(`http://localhost:8080/api/document-requests/${id}`, {
      headers: {
        Cookie: `auth_cookie=${token}`,
        "Content-Type": "application/json",
      },
      cache: "no-store",
    });

    if (res.status === 403) {
      return { data: null, error: "FORBIDDEN" };
    }

    if (res.status === 404) {
      return { data: null, error: "NOT_FOUND" };
    }

    if (!res.ok) {
      return { data: null, error: "SERVER_ERROR" };
    }

    const body: ApiResponse<DocumentRequest> = await res.json();

    return {
      data: body.data || null,
      error: body.data ? undefined : "NOT_FOUND",
    };
  } catch (error) {
    console.error("Fetch error:", error);
    return { data: null, error: "SERVER_ERROR" };
  }
}
