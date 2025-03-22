FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY cmd/ cmd/
COPY internal/ internal/

# Build the Go app
RUN go build -o go_block_go ./cmd/main.go

FROM ubuntu:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder --chown=ubuntu:ubuntu /app/go_block_go .

# non-root user
USER 1000

ENTRYPOINT ["./go_block_go"]