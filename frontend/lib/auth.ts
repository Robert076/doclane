import { cookies } from "next/headers";
import { User, ApiResponse } from "@/types";

export async function getCurrentUser(): Promise<User | null> {
  const cookieStore = await cookies();
  const token = cookieStore.get("auth_cookie")?.value;

  if (!token) return null;

  try {
    const res = await fetch("http://localhost:8080/api/users/me", {
      headers: {
        Cookie: `auth_cookie=${token}`,
        "Content-Type": "application/json",
      },
      cache: "no-store",
    });

    if (!res.ok) {
      if (res.status !== 401) {
        console.error("Go API responded with:", res.status);
      }
      return null;
    }

    const body: ApiResponse<User> = await res.json();
    return body.data;
  } catch (error) {
    return null;
  }
}
