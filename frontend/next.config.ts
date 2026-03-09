import type { NextConfig } from "next";

const nextConfig: NextConfig = {
        async rewrites() {
                return [
                        {
                                source: "/api/backend/:path*",
                                destination:
                                        "https://8mcdcahxsi.execute-api.eu-west-1.amazonaws.com/api/:path*",
                        },
                ];
        },
};

export default nextConfig;
