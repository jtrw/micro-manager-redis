# Micro Manager Redis

[![Build](https://github.com/jtrw/micro-manager-redis/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/jtrw/micro-manager-redis/actions)
[![codecov](https://codecov.io/gh/jtrw/micro-manager-redis/graph/badge.svg?token=R4WZPK20B7)](https://codecov.io/gh/jtrw/micro-manager-redis)

Micro Manager Redis is a simple tool for managing keys in a Redis database.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Endpoints](#endpoints)
    - [API Endpoints](#api-endpoints)
    - [Authentication](#authentication)
    - [File Server](#file-server)
- [Testing](#testing)
- [License](#license)

## Installation

To install Micro Manager Redis, follow these steps for backend:
```bash
# Clone the repository
git clone https://github.com/your-username/micro-manager-redis.git

# Change into the project directory
cd micro-manager-redis/backend

# Build the application
go build -o micro-manager-redis main.go

# Run the application
./micro-manager-redis
```

For installation of the web UI, follow these steps:
```bash
cd frontend
yarn dev
```

## Usage

Micro Manager Redis provides a set of HTTP APIs for managing keys in a Redis database. The default server address is http://localhost:8080.

docker-compose.yml
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

## Configuration

The tool can be configured using command-line options or environment variables. Available configuration options:

* `-l, --listen`: Listen address for the server (default: :8080).
* `-s, --secret`: Secret key for secure operations (default: 123).
* `--pinsize`: Size of PIN for secure operations (default: 5).
* `--expire`: Maximum lifetime for keys (default: 24h).
* `--pinattempts`: Maximum attempts to enter PIN (default: 3).
* `--web`: Web UI location (default: ./web).
* `--redis-url`: Redis server URL (default: localhost:6379).
* `--redis-db`: Redis database name (default: 3).
* `--redis-pass`: Redis database password (default: Y6zhcj769Fo1).
* `--auth-login`: Authentication login (default: admin).
* `--auth-password`: Authentication password (default: admin).

## Endpoints

### API Endpoints
Micro Manager Redis exposes the following API endpoints:

* `GET /api/v1/keys`: Retrieve all keys.
* `GET /api/v1/keys/{key}`: Retrieve details of a specific key.
* `DELETE /api/v1/keys/{key}`: Delete a specific key.
* `DELETE /api/v1/keys`: Delete all keys.
* `GET /api/v1/keys-group`: Group keys based on a separator.
* `DELETE /api/v1/keys-group/{group}`: Delete keys based on a group.

### Authentication

Authentication is required for API endpoints. Use the Bearer token in the `Authorization` header.

### File Server

Micro Manager Redis serves static files from the /web directory.

## Testing

To run tests, execute the following command:
```bash
go test -v ./...
```

## License

Micro Manager Redis is open-source software licensed under the [MIT license](https://opensource.org/licenses/MIT).
```
