FROM golang:1.23-bookworm AS builder

WORKDIR /go/src/github.com/forbole/bdjuno

ENV CGO_CFLAGS="-O -D__BLST_PORTABLE__"
ENV CGO_CFLAGS_ALLOW="-O -D__BLST_PORTABLE__"

ENV GOPRIVATE=github.com/mocachain/*
ENV GOPROXY=https://proxy.golang.org,direct
ENV GONOSUMDB=github.com/mocachain/*
ENV GONOSUMCHECK=github.com/mocachain/*

ARG GITHUB_TOKEN
RUN if [ -n "$GITHUB_TOKEN" ]; then \
        git config --global url."https://${GITHUB_TOKEN}:@github.com/".insteadOf "https://github.com/"; \
    fi

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make build

ADD https://github.com/hasura/graphql-engine/releases/download/v2.43.0/cli-hasura-linux-amd64 ./build/hasura
RUN chmod +x ./build/hasura


FROM golang:1.23-bookworm

WORKDIR /root

COPY --from=builder /go/src/github.com/forbole/bdjuno/build/hasura /usr/bin/hasura
COPY --from=builder /go/src/github.com/forbole/bdjuno/hasura ./hasura
COPY --from=builder /go/src/github.com/forbole/bdjuno/build/bdjuno /usr/bin/bdjuno

CMD [ "bdjuno" ]