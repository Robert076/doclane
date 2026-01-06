// lib/getFilesByRequestId.ts
import { cookies } from "next/headers";

export default async function getFilesByRequestId(requestId: string) {
  // 1. Așteptăm rezolvarea cookies (Next.js 15+)
  const cookieStore = await cookies();

  // 2. Extragem string-ul de cookies pentru header
  const allCookies = cookieStore.toString();

  try {
    const res = await fetch(`http://localhost:8080/api/document-requests/${requestId}/files`, {
      method: "GET",
      headers: {
        // Trimitem toate cookie-urile (inclusiv token-ul JWT) către Go
        Cookie: allCookies,
        "Content-Type": "application/json",
      },
      // Important: URL-urile presemnate se schimbă, deci nu vrem cache agresiv
      cache: "no-store",
    });

    if (!res.ok) {
      console.error("Fetch files failed:", res.status);
      return null;
    }

    return await res.json();
  } catch (error) {
    console.error("Error fetching files:", error);
    return null;
  }
}
