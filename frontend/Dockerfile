# --------------------------------------------------
# Production stage
FROM node:20-alpine AS deploy
WORKDIR /app
ENV NODE_ENV=production
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --chown=nextjs:nodejs .next/standalone ./.next/standalone
COPY --chown=nextjs:nodejs .next/static ./.next/static
COPY --chown=nextjs:nodejs public ./public

USER nextjs

EXPOSE 8080
ENV PORT=8080
ENV HOSTNAME="0.0.0.0"

# デバッグ用: 起動時にディレクトリ構造を確認
CMD ["sh", "-c", "echo '=== Checking directory structure ===' && \
    ls -la /app && \
    echo '=== Checking .next directory ===' && \
    ls -la /app/.next && \
    echo '=== Checking .next/standalone ===' && \
    ls -la /app/.next/standalone && \
    echo '=== Checking .next/static ===' && \
    ls -la /app/.next/static && \
    echo '=== Starting Next.js server ===' && \
    node /app/.next/standalone/server.js"]

# --------------------------------------------------
# dev
FROM node:20-alpine AS dev
LABEL build-stage=installer
RUN apk add --no-cache libc6-compat
RUN npm install -g pnpm
WORKDIR /app
COPY package.json ./
COPY pnpm-lock.yaml ./
COPY . .
ENV SKIP_PREPARE='true'
RUN pnpm install --frozen-lockfile
EXPOSE 3000
CMD ["pnpm", "run", "dev"]
