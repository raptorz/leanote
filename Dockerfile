# Stage 1: Build the application
FROM golang:1.24-alpine AS builder

# Set Go proxy for China
ENV GOPROXY=https://goproxy.cn,direct

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate the Revel application entrypoint, then build the server rather than
# the Revel command-line helper in app/cmd.
RUN go run ./app/cmd build . /tmp/pearlnote-revel-build prod && \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -o /pearlnote -ldflags="-s -w" ./app/tmp

RUN REVEL_DIR="$(go list -m -f '{{.Dir}}' github.com/revel/revel)" && \
    mkdir -p /runtime/github.com/revel/revel && \
    cp -R "$REVEL_DIR/conf" "$REVEL_DIR/templates" /runtime/github.com/revel/revel/

# Stage 2: Runtime image
FROM alpine:3.22

LABEL maintainer="raptor<raptor.zh@gmail.com>"

# Install ca-certificates for HTTPS connections and timezone data
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /opt/pearlnote

# Copy the binary from builder
COPY --from=builder /pearlnote /opt/pearlnote/pearlnote
COPY --from=builder /runtime /opt/pearlnote/runtime

# Copy necessary runtime files
COPY conf/ /opt/pearlnote/conf/
COPY messages/ /opt/pearlnote/messages/
COPY public/ /opt/pearlnote/public/
COPY app/views/ /opt/pearlnote/app/views/
RUN mkdir -p /opt/pearlnote/runtime/github.com/pearlnote && \
    ln -s /opt/pearlnote /opt/pearlnote/runtime/github.com/pearlnote/pearlnote

EXPOSE 9000

# Set working directory to where the app needs to run
WORKDIR /opt/pearlnote

CMD ["/opt/pearlnote/pearlnote", "-importPath", "github.com/pearlnote/pearlnote", "-srcPath", "/opt/pearlnote/runtime", "-runMode", "prod"]
