# micro-manager-caches-keys

[![Build](https://github.com/jtrw/micro-manager-redis/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/jtrw/micro-manager-redis/actions)
[![codecov](https://codecov.io/gh/jtrw/micro-manager-redis/graph/badge.svg?token=R4WZPK20B7)](https://codecov.io/gh/jtrw/micro-manager-redis)

## Installation

```bash
version: "3"
services:
    manage-rkeys:
    container_name: manage-rkeys
    build:
        context: .
        dockerfile: Dockerfile
        args:
        MANAGE_RKEYS_URL: ${MANAGE_RKEYS_URL}
    ports:
        - "8080:8080"
    extra_hosts:
        - "host.docker.internal:host-gateway"
    env_file:
        - .env
    environment:
        - REDIS_URL=${REDIS_URL}
        - REDIS_DATABASE=${REDIS_DATABASE}
        - REDIS_PASSWORD=${REDIS_PASSWORD}
        - AUTH_LOGIN=${AUTH_LOGIN}
        - AUTH_PASSWORD=${AUTH_PASSWORD}
        - LISTEN_SERVER=${LISTEN_SERVER}
    healthcheck:
        test:
        [
            "CMD",
            "sh",
            "-c",
            "wget -qO- http://127.0.0.1:8080/ping | grep -q 'pong' || exit 1",
        ]
        interval: 5s
        timeout: 3s
        retries: 3
```


Nginx configuration
```
location /index.html {
         proxy_redirect          off;
         proxy_set_header        Host $http_host;
         proxy_pass              http://127.0.1:8080/web/index.html;

```
