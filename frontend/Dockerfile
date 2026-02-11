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
RUN pnpm run build

# --------------------------------------------------
# Production stage
FROM node:18-alpine AS deploy
WORKDIR /app
ENV NODE_ENV=production
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs
COPY --from=deps --chown=nextjs:nodejs /app/node_modules ./node_modules
COPY --from=deps --chown=nextjs:nodejs /app/package.json ./package.json
COPY --from=deps --chown=nextjs:nodejs /app/.next ./.next
COPY --chown=nextjs:nodejs ./public ./public
COPY --chown=nextjs:nodejs ./next.config.ts ./

USER nextjs

EXPOSE 8080

ENV PORT=8080
ENV NO_NEXT_TRANSPILE='true'

CMD ["npm", "run", "start"]
