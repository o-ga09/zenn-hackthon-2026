# Build stage
FROM node:18-alpine AS deps
LABEL build-stage=installer
RUN apk add --no-cache libc6-compat
RUN npm install -g pnpm
WORKDIR /app
COPY package.json ./
COPY pnpm-lock.yaml ./
ENV SKIP_PREPARE='true'
RUN pnpm install --prod --frozen-lockfile

# --------------------------------------------------
# Production stage
FROM node:18-alpine AS deploy
WORKDIR /app
ENV NODE_ENV=production
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Copy standalone build (includes minimal node_modules)
COPY --chown=nextjs:nodejs ./.next/standalone ./
# Copy static assets
COPY --chown=nextjs:nodejs ./.next/static ./.next/static
# Copy public directory
COPY --chown=nextjs:nodejs ./public ./public

USER nextjs

EXPOSE 8080

ENV PORT=8080
ENV HOSTNAME="0.0.0.0"

CMD ["node", "server.js"]
