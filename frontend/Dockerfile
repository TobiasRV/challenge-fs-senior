FROM node:25-alpine3.22 AS builder
WORKDIR /usr/src/app
COPY ./package.json ./package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build && \
    npm prune --production

FROM node:25-alpine3.22
ENV NODE_ENV=production
RUN apk update && apk upgrade --no-cache
WORKDIR /usr/src/app
COPY --chown=node:node --from=builder /usr/src/app ./
USER node
ENTRYPOINT ["npm", "run", "start"]
