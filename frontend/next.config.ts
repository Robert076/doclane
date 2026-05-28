import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // "standalone" emits .next/standalone with a minimal node_modules tree
  // so the production image can run `node server.js` without the full repo.
  // Required for our Docker runtime stage to stay small.
  output: "standalone",
};

export default nextConfig;
