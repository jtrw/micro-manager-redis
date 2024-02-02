FROM --platform=$BUILDPLATFORM node:21.4.0-alpine AS build-frontend

ARG SKIP_FRONTEND_TEST
ARG SKIP_FRONTEND_BUILD

ARG ARG_MANAGE_RKEYS_URL
ENV MANAGE_RKEYS_URL=${ARG_MANAGE_RKEYS_URL}

WORKDIR /srv/frontend/

COPY ./frontend/ /srv/frontend/
RUN echo "SKIP_FRONTEND_TEST=$SKIP_FRONTEND_TEST"
RUN echo "SKIP_FRONTEND_BUILD=$SKIP_FRONTEND_BUILD"
RUN apk add --no-cache --update git && \
    yarn install --frozen-lockfile

RUN yarn build

FROM golang:1.21.5-alpine AS build-backend

ENV GOFLAGS="-mod=vendor"
ENV CGO_ENABLED=0

ARG CI
ARG GITHUB_REF
ARG GITHUB_SHA
ARG GIT_BRANCH
ARG SKIP_BACKEND_TEST
ARG BACKEND_TEST_TIMEOUT
ARG GIT_BRANCH
ARG GITHUB_SHA

ADD backend /build/backend
COPY --from=build-frontend /srv/frontend/dist/ /build/backend/web/
#RUN find /build/backend/web/ -regex '.*\.\(html\|js\|mjs\)$' -print -exec sed -i "s|{% RKEYS_URL %}|http://127.0.0.1:8080|g" {} \;
WORKDIR /build/backend

# install gcc in order to be able to go test package with -race
RUN apk --no-cache add gcc libc-dev

RUN apk add --no-cache --update git tzdata ca-certificates bash

RUN go mod vendor

RUN \
    if [ -z "$CI" ] ; then \
    echo "runs outside of CI" && version=$(git rev-parse --abbrev-ref HEAD)-$(git log -1 --format=%h)-$(date +%Y%m%dT%H:%M:%S); \
    else version=${GIT_BRANCH}-${GITHUB_SHA:0:7}-$(date +%Y%m%dT%H:%M:%S); fi && \
    echo "version=$version"

RUN echo go version: `go version`

# run tests
#RUN \
#    cd app && \
#    if [ -z "$SKIP_BACKEND_TEST" ] ; then \
#        CGO_ENABLED=1 go test -race -p 1 -timeout="${BACKEND_TEST_TIMEOUT:-300s}" -covermode=atomic \
#    else \
#        echo "Skip tests" \
#    ; fi


RUN cd app && go build -o rkeys -ldflags "-X main.revision=${version} -s -w"

#FROM scratch
FROM alpine

ARG ARG_MANAGE_RKEYS_URL
ENV MANAGE_RKEYS_URL=${ARG_MANAGE_RKEYS_URL}
ARG GITHUB_SHA

RUN echo "prepare environment"
# replace {% MANAGE_RKEYS_URL %} by content of MANAGE_RKEYS_URL variable
RUN find . -regex '.*\.\(html\|js\|mjs\)$' -print -exec sed -i "s|{% MANAGE_RKEYS_URL %}|${MANAGE_RKEYS_URL}|g" {} \;

LABEL org.opencontainers.image.authors="Nil Borodulia <nil.borodulia@gmail.com>" \
      org.opencontainers.image.description="Manager contents redis keys" \
      org.opencontainers.image.documentation="https://github.com/jtrw/micro-manager-redis" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.source="https://github.com/jtrw/micro-manager-redis.git" \
      org.opencontainers.image.title="ManagerRedis" \
      org.opencontainers.image.revision="${GITHUB_SHA}"

WORKDIR /srv

COPY docker-init.sh /srv/init.sh
RUN chmod +x /srv/init.sh

COPY --from=build-backend /build/backend/app/rkeys /srv/rkeys
COPY --from=build-frontend /srv/frontend/dist/ /srv/web/


#RUN chown -R app:app /srv
#RUN ln -s /srv/rkeys /usr/bin/rkeys

#EXPOSE 8080
#CMD ["/srv/rkeys"]
ENTRYPOINT ["/srv/init.sh"]
