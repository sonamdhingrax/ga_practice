# syntax=docker/dockerfile:1

# ---- Build stage ---------------------------------------------------------
# Pinned Go version keeps local + CI builds reproducible.
FROM --platform=$BUILDPLATFORM golang:1.25 AS build

# TARGETOS/TARGETARCH are provided automatically by buildx for each platform
# in the --platform list, enabling cross-compilation from a single builder.
ARG TARGETOS
ARG TARGETARCH

# Build metadata injected by CI (defaults keep `docker build` working locally).
ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown

WORKDIR /src

# Copy go.mod first so dependency download is cached independently of source
# changes (classic Docker layer-cache optimization).
COPY go.mod ./
# COPY go.sum ./   # uncomment once you add external dependencies
RUN go mod download

COPY . .

# CGO disabled => a fully static binary that runs on distroless/static.
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath \
    -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" \
    -o /out/app .

# ---- Runtime stage -------------------------------------------------------
# distroless/static: no shell, no package manager, runs as nonroot by default.
FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY --from=build /out/app /app

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app"]
