FROM --platform=$BUILDPLATFORM node:21.4.0-alpine AS build-frontend

ARG SKIP_FRONTEND_TEST
ARG SKIP_FRONTEND_BUILD

ARG MANAGE_RKEYS_URL
ENV MANAGE_RKEYS_URL=$MANAGE_RKEYS_URL

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

ARG GITHUB_SHA

WORKDIR /srv

#COPY docker-init.sh /srv/init.sh

COPY --from=build-backend /build/backend/app/rkeys /srv/rkeys
COPY --from=build-frontend /srv/frontend/dist/ /srv/web/
#RUN chown -R app:app /srv
#RUN ln -s /srv/rkeys /usr/bin/rkeys

#EXPOSE 8080
#HEALTHCHECK --interval=30s --timeout=3s CMD curl --fail http://localhost:8080/ping || exit 1
#HEALTHCHECK --interval=5s --timeout=3s CMD wget -qO- http://127.0.0.1:8080/ping | grep -q 'pong' || exit 1
CMD ["/srv/rkeys", "server"]
#ENTRYPOINT ["/srv/rkeys", "server"]
