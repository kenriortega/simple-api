# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.16-alpine as builder

# Add Maintainer Info
LABEL maintainer="Kenrique Ortega <kenriortega@gmail.com>"

ENV GO111MODULE=on

RUN apk --no-cache add ca-certificates build-base upx

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Build the Go app
RUN go build -ldflags "-s -w"  -o /app/api main.go
RUN upx -9 /app/api



######## Start a new stage from scratch #######
FROM alpine:latest

WORKDIR /app/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/api .

EXPOSE 3000

ENV DATABASE_URL=mongodb://localhost:27017/db

ENV PORT=3000
ENTRYPOINT ["/app/api"]
