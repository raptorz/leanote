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

# Build the application directly using go build
# The main entry point is in app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /pearlnote -ldflags="-s -w" github.com/pearlnote/pearlnote/app/cmd

# Stage 2: Runtime image
FROM alpine:3.22

LABEL maintainer="raptor<raptor.zh@gmail.com>"

# Install ca-certificates for HTTPS connections and timezone data
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /opt/pearlnote

# Copy the binary from builder
COPY --from=builder /pearlnote /opt/pearlnote/pearlnote

# Copy necessary runtime files
COPY conf/ /opt/pearlnote/conf/
COPY messages/ /opt/pearlnote/messages/
COPY public/ /opt/pearlnote/public/
COPY app/views/ /opt/pearlnote/app/views/

EXPOSE 9000

# Set working directory to where the app needs to run
WORKDIR /opt/pearlnote

CMD ["/opt/pearlnote/pearlnote", "run", "github.com/pearlnote/pearlnote", "-importPath", "github.com/pearlnote/pearlnote"]