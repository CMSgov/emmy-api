FROM node:25.8.1-slim AS base

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

RUN npm install -g --force corepack
RUN corepack enable pnpm

WORKDIR /app

FROM base AS builder

COPY package.json package-lock.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
COPY . .
RUN pnpm build

FROM base

COPY --from=builder --chown=node:node /app/dist/index.js .

EXPOSE 3000
CMD ["index.js"]
