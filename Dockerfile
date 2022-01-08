# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.16-alpine as builder

# Add Maintainer Info
LABEL maintainer="Kenrique Ortega <kenriortega@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

RUN apk --no-cache add ca-certificates build-base musl-dev

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -ldflags "-s -w"  -o /app/api main.go



######## Start a new stage from scratch #######
FROM alpine:latest

# RUN apk --no-cache add ca-certificates
# RUN apk add librdkafka-dev pkgconf

# Build Args
ARG APP_DIR=/app

# Create APP Directory
RUN mkdir -p ${APP_DIR}
WORKDIR /app/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/api .

EXPOSE 8080

ENV DATABASE_URL=mongodb://localhost:27017/db
ENV PORT=8080
ENTRYPOINT ["/app/api"]
