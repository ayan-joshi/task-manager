import { fileURLToPath } from 'node:url';
import { dirname } from 'node:path';

const __dirname = dirname(fileURLToPath(import.meta.url));

/** @type {import('next').NextConfig} */
const nextConfig = {
  // Produces a minimal standalone server bundle for small Docker images.
  output: 'standalone',
  // Pin file-tracing to this app so the standalone build is correct even when
  // other lockfiles exist higher in the directory tree.
  outputFileTracingRoot: __dirname,
  reactStrictMode: true,
  poweredByHeader: false,
};

export default nextConfig;
