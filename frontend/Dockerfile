FROM oven/bun:1 as frontend

# Define build arguments for environment variables
# These args will be available as env vars
ARG VITE_API_LANGUAGE=fa
ARG VITE_APP_ENVIRONMENT
ARG VITE_AUTH_TOKEN
ARG VITE_VARIANT=cab
ARG VITE_CENTRAL_BACKEND_API
ARG VITE_ENV

RUN apt-get update && apt-get install -y git

RUN mkdir -p /app/
WORKDIR /app

COPY package.json bun.lockb ./

RUN bun i --frozen-lockfile

COPY . .

RUN bun run build

FROM nginxinc/nginx-unprivileged

COPY --from=frontend ./app/dist/ /usr/share/nginx/html

USER root
RUN chgrp -R 0 /usr/share/nginx && chmod -R g=u /usr/share/nginx

USER nginx
COPY .docker/EntryPoint.sh /EntryPoint.sh

ENTRYPOINT ["/EntryPoint.sh"]
