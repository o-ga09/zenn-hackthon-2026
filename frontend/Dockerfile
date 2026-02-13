# --------------------------------------------------
# Production stage
FROM node:20-alpine AS deploy
WORKDIR /app
ENV NODE_ENV=production
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --chown=nextjs:nodejs ./.next/standalone ./
COPY --chown=nextjs:nodejs ./.next/static ./.next/static
COPY --chown=nextjs:nodejs ./public ./public

USER nextjs

EXPOSE 8080
ENV PORT=8080
ENV HOSTNAME="0.0.0.0"

CMD ["node", "server.js"]

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
