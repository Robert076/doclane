import { cookies } from "next/headers";
import { NextRequest, NextResponse } from "next/server";

export async function POST(req: NextRequest) {
  const { idToken, accessToken } = await req.json();

  if (!idToken || !accessToken) {
    return NextResponse.json({ error: "Missing tokens" }, { status: 400 });
  }

  const cookieStore = await cookies();
  cookieStore.set("auth_cookie", idToken, {
    httpOnly: true,
    sameSite: "lax",
    path: "/",
    expires: new Date(Date.now() + 1000 * 60 * 60),
  });
  cookieStore.set("access_token", accessToken, {
    httpOnly: true,
    sameSite: "lax",
    path: "/",
    expires: new Date(Date.now() + 1000 * 60 * 60),
  });

  return NextResponse.json({ success: true });
}
