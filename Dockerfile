## How to build: docker build --build-arg GITHUB_TOKEN=$GITHUB_TOKEN -t tag .

# FROM alpine:3.18
# RUN apk add -y --no-cache libc6-compat
# FROM debian:12.2-slim
FROM golang:1.25.0-bookworm AS builder

# Install git
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        git \
        curl \
        xz-utils && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install upx for compressing Go binaries
RUN curl -L https://github.com/upx/upx/releases/download/v5.0.2/upx-5.0.2-amd64_linux.tar.xz \
    -o /tmp/upx.tar.xz && \
    mkdir -p /tmp/upx && \
    tar -xf /tmp/upx.tar.xz -C /tmp/upx --strip-components=1 && \
    mv /tmp/upx/upx /usr/local/bin/upx && \
    chmod +x /usr/local/bin/upx && \
    rm -rf /tmp/upx*

# Install orchestrion execute file for Datadog APM. 
RUN go install github.com/DataDog/orchestrion@v1.5.0

# Set the Current Working Directory inside the container
WORKDIR /build

# Copy all files to the container
COPY . .

RUN echo "show workdir before build" && \
    pwd && \
    ls -al ./

# Setup GITHUB_TOKEN
ARG GITHUB_TOKEN
ENV GITHUB_TOKEN=${GITHUB_TOKEN}

# git auth use PAT
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

# Setup /build to be a safe directory if there are git commands in Dockerfile
RUN git config --global --add safe.directory /build

# Set private module access
RUN go env -w GOPRIVATE=github.com/phcp-tech

# Download dependencies
RUN go mod download

# Build the Go app
RUN GOEXPERIMENT=jsonv2,greenteagc CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-s -w"

# Compress the Go app
RUN upx --best --lzma template

RUN echo "show workdir after build" && \
    pwd && \
    ls -al ./

# Start a new stage from scratch. Use ubuntu:22.04 because there are no ps/top commands in debian:12.2-slim
# FROM debian:12.2-slim
FROM ubuntu:22.04

# Install timezone data, curl and certificates
RUN set -e && \
    apt-get update --fix-missing && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
        tzdata \
        curl \
        vim \
        ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Set the current working directory inside the container
ENV APP_HOME=/opt/phcp-tech/template
WORKDIR $APP_HOME 

# Copy executable file
RUN mkdir -p $APP_HOME $APP_HOME/config $APP_HOME/logs
COPY --from=builder /build/template $APP_HOME/template
COPY --from=builder /build/config $APP_HOME/config
COPY --from=builder /build/logs $APP_HOME/logs

# Soft link for logs to download
RUN ln -s /var/log/template.log $APP_HOME/logs/template.log 

# Expose port
EXPOSE 80

# Run the executable
# CMD ./template > /var/log/template.log 2>&1
CMD ["./template"]
