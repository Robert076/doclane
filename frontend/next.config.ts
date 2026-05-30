import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  // Trust X-Forwarded-* headers from the ALB so cookies and redirects
  // work correctly when TLS is terminated at the load balancer.
  serverExternalPackages: [],
  allowedDevOrigins: [],
};

export default nextConfig;
