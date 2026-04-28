import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

const PUBLIC_PATHS = ["/login", "/register"];

export default function middleware(request: NextRequest) {
        const token = request.cookies.get("auth_cookie")?.value;
        const { pathname } = request.nextUrl;

        const isPublic = PUBLIC_PATHS.some((p) => pathname.startsWith(p));

        if (!token && !isPublic) {
                const url = request.nextUrl.clone();
                url.pathname = "/login";
                url.searchParams.set("callbackUrl", pathname);
                return NextResponse.redirect(url);
        }

        if (token && isPublic) {
                const url = request.nextUrl.clone();
                url.pathname = "/dashboard/requests";
                return NextResponse.redirect(url);
        }

        if (token && pathname === "/") {
                const url = request.nextUrl.clone();
                url.pathname = "/dashboard/requests";
                return NextResponse.redirect(url);
        }

        return NextResponse.next();
}

export const config = {
        matcher: ["/((?!api|_next/static|_next/image|favicon.ico|.*\\..*).*)"],
};
