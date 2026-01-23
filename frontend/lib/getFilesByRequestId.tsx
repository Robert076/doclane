import { ApiResponse, DocumentFile } from "@/types";
import { cookies } from "next/headers";

export default async function getFilesByRequestId(
  requestId: string
): Promise<ApiResponse<DocumentFile[]> | null> {
  const cookieStore = await cookies();
  const allCookies = cookieStore.toString();

  try {
    const res = await fetch(`http://localhost:8080/api/document-requests/${requestId}/files`, {
      method: "GET",
      headers: {
        Cookie: allCookies,
        "Content-Type": "application/json",
      },
      cache: "no-store",
    });

    if (!res.ok) {
      console.error(`Fetch files failed with status: ${res.status}`);
      return null;
    }

    const result: ApiResponse<DocumentFile[]> = await res.json();

    return result;
  } catch (error) {
    console.error("Error fetching files:", error);
    return null;
  }
}
