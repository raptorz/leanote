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
RUN CGO_ENABLED=0 GOOS=linux go build -o /leanote -ldflags="-s -w" github.com/leanote/leanote/app/cmd

# Stage 2: Runtime image
FROM alpine:3.22

LABEL maintainer="raptor<raptor.zh@gmail.com>"

# Install ca-certificates for HTTPS connections and timezone data
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /opt/leanote

# Copy the binary from builder
COPY --from=builder /leanote /opt/leanote/leanote

# Copy necessary runtime files
COPY conf/ /opt/leanote/conf/
COPY messages/ /opt/leanote/messages/
COPY public/ /opt/leanote/public/
COPY app/views/ /opt/leanote/app/views/

EXPOSE 9000

# Set working directory to where the app needs to run
WORKDIR /opt/leanote

CMD ["/opt/leanote/leanote", "run", "github.com/leanote/leanote", "-importPath", "github.com/leanote/leanote"]