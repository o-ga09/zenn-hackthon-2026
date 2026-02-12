import type { NextConfig } from 'next'

const nextConfig: NextConfig = {
  output: 'standalone',
  assetPrefix: process.env.NODE_ENV === 'production' ? 'https://static-agent.tavinikkiy.com' : '',
  /* config options here */
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: '**.r2.cloudflarestorage.com',
        port: '',
        pathname: '/**',
      },
      {
        protocol: 'https',
        hostname: '**.r2.dev',
        port: '',
        pathname: '/**',
      },
      {
        protocol: 'https',
        hostname: 'lh3.googleusercontent.com',
        port: '',
        pathname: '/**',
      },
      {
        protocol: 'https',
        hostname: 'firebasestorage.googleapis.com',
        port: '',
        pathname: '/**',
      },
      {
        protocol: 'https',
        hostname: 'example.com',
        port: '',
        pathname: '/**',
      },
    ],
  },
}

export default nextConfig
