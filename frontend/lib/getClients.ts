import { User, ApiResponse } from "@/types";
import { cookies } from "next/headers";

export default async function getMyClients(): Promise<User[]> {
  const cookieStore = await cookies();
  const token = cookieStore.get("auth_cookie")?.value;

  try {
    const res = await fetch("http://localhost:8080/api/users/my-clients", {
      headers: {
        Cookie: `auth_cookie=${token}`,
      },
      cache: "no-store",
    });

    if (!res.ok) return [];

    const body: ApiResponse<User[]> = await res.json();
    return body.data || [];
  } catch (error) {
    console.error("Fetch clients error:", error);
    return [];
  }
}
