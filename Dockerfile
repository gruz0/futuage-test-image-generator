# =============================================================================
# FutuAge Test Image Generator - Multi-stage Docker Build
# =============================================================================
# Uses Alpine for building with static linking, scratch for minimal runtime
# Resulting image is ~15-20MB with full WebP support via CGO
# =============================================================================

# -----------------------------------------------------------------------------
# Build Stage
# -----------------------------------------------------------------------------
FROM golang:1.25-alpine AS builder

# Install build dependencies for CGO and static linking
# - gcc, musl-dev: C compiler and standard library
# - libwebp-dev, libwebp-static: WebP library (dynamic and static)
RUN apk add --no-cache \
    gcc \
    musl-dev \
    libwebp-dev \
    libwebp-static

WORKDIR /build

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with CGO enabled and static linking for scratch compatibility
# -s: Omit symbol table
# -w: Omit DWARF debugging info
# -linkmode external: Use external linker
# -extldflags '-static': Build fully static binary
ENV CGO_ENABLED=1
RUN go build \
    -ldflags="-s -w -linkmode external -extldflags '-static'" \
    -o /futuage-test-image-gen \
    .

# -----------------------------------------------------------------------------
# Runtime Stage (scratch = empty image)
# -----------------------------------------------------------------------------
FROM scratch

# Copy the statically linked binary
COPY --from=builder /futuage-test-image-gen /futuage-test-image-gen

# Copy default configs (can be overwritten via volume mount)
COPY --from=builder /build/configs /configs

# Create output directory structure marker
# Note: scratch has no shell, directories are created by the app
WORKDIR /output

# Default command generates all images to /output
ENTRYPOINT ["/futuage-test-image-gen"]
CMD ["generate", "--output", "/output"]

