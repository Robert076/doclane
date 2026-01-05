import { ApiResponse, DocumentRequest, UserRole } from "@/types";
import { cookies } from "next/headers";

export default async function getDocumentRequests(role: UserRole): Promise<DocumentRequest[]> {
  const cookieStore = await cookies();
  const token = cookieStore.get("auth_cookie")?.value;

  try {
    const res = await fetch(
      `http://localhost:8080/api/document-requests/${role.toLowerCase()}/documents`,
      {
        headers: {
          Cookie: `auth_cookie=${token}`,
          "Content-Type": "application/json",
        },
        cache: "no-store",
      }
    );

    if (!res.ok) {
      console.error(`Failed to fetch document requests for ${role}`);
      return [];
    }

    const body: ApiResponse<DocumentRequest[]> = await res.json();
    return body.data || [];
  } catch (error) {
    console.error("Fetch error:", error);
    return [];
  }
}
