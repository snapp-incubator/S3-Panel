FROM node:22-alpine AS frontend

# Build-time configuration, exposed to Vite as environment variables.
ARG VITE_API_LANGUAGE=fa
ARG VITE_APP_ENVIRONMENT
ARG VITE_AUTH_TOKEN
ARG VITE_VARIANT=cab
ARG VITE_CENTRAL_BACKEND_API
ARG VITE_ENV
ENV VITE_API_LANGUAGE=$VITE_API_LANGUAGE \
    VITE_APP_ENVIRONMENT=$VITE_APP_ENVIRONMENT \
    VITE_AUTH_TOKEN=$VITE_AUTH_TOKEN \
    VITE_VARIANT=$VITE_VARIANT \
    VITE_CENTRAL_BACKEND_API=$VITE_CENTRAL_BACKEND_API \
    VITE_ENV=$VITE_ENV

RUN corepack enable && apk add --no-cache git

WORKDIR /app

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile

COPY . .
RUN pnpm run build

FROM nginxinc/nginx-unprivileged

COPY --from=frontend /app/dist/ /usr/share/nginx/html
COPY .docker/nginx.conf /etc/nginx/conf.d/default.conf

USER root
RUN chgrp -R 0 /usr/share/nginx && chmod -R g=u /usr/share/nginx

USER nginx
COPY .docker/EntryPoint.sh /EntryPoint.sh

ENTRYPOINT ["/EntryPoint.sh"]
