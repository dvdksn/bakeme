# syntax=docker/dockerfile:1

ARG GO_VERSION="1.23"
ARG GOLANGCI_LINT_VERSION="1.61"

# base downloads the necessary Go modules
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS base
WORKDIR /src

RUN --mount=src=go.mod,dst=go.mod \
	--mount=src=go.sum,dst=go.sum \
	--mount=type=cache,target=/go/pkg/mod \
	go mod download

# build compiles the program
FROM base AS build
ARG TARGETOS TARGETARCH BUILD_TAGS
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
RUN --mount=target=. \
	--mount=type=cache,target=/go/pkg/mod \
	go build -tags="${BUILD_TAGS}" -o "/usr/bin/bakeme" .

# image creates a runtime image
FROM alpine AS image
COPY --from=build "/usr/bin/bakeme" "/usr/bin/bakeme"
ENTRYPOINT ["/usr/bin/bakeme"]

FROM base AS test
RUN --mount=target=. \
	--mount=type=cache,target=/go/pkg/mod \
	go test .

FROM golangci/golangci-lint:v${GOLANGCI_LINT_VERSION}-alpine AS lint
RUN --mount=target=.,rw \
	golangci-lint run
