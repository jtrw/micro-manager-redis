
OS=linux
ARCH=amd64

dockerx:
    docker buildx build --progress=plain --platform linux/amd64,linux/arm/v7,linux/arm64 --no-cache -t jtrw/micro-manager-redis:latest --push .
