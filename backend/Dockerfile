# syntax=docker/dockerfile:1.7

# ---------- Build stage ----------
# Pinned to the Go version declared in backend/go.mod.
FROM golang:1.25-alpine AS build

# git is needed for some Go module fetches; ca-certificates for HTTPS to proxy.golang.org.
RUN apk add --no-cache git ca-certificates

WORKDIR /src

# Copy module files first so dependency download is cached independently of source changes.
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Now copy the rest of the backend source.
COPY . .

# Build a fully static binary so it can run on distroless/static.
# CGO_ENABLED=0 is critical: the lib/pq driver is pure Go, so no C deps are needed,
# and this lets us use a scratch-based runtime.
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ENV CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH}

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /out/doclane-backend ./cmd

# ---------- Runtime stage ----------
# distroless/static is the smallest secure base for static Go binaries.
# It includes ca-certificates, /etc/passwd with a nonroot user, and tzdata. No shell.
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=build /out/doclane-backend /app/doclane-backend

# Non-root by default (uid 65532 from the :nonroot tag).
USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/app/doclane-backend"]
