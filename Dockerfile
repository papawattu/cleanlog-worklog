# Use a multi-stage build to support multiple architectures
# Stage 1: Build stage
FROM golang:1.23.1-alpine AS builder
LABEL org.opencontainers.image.source=https://github.com/papawattu/cleanlog-worklog
LABEL org.opencontainers.image.description="A simple web app log cleaning house"
LABEL org.opencontainers.image.licenses=MIT

ARG USER=nouser

RUN apk add --no-cache make

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build


# Stage 2: Final stage
FROM alpine AS build-stage

ARG USER=nouser

WORKDIR /

COPY --from=builder /app/bin/worklog /worklog

RUN adduser -D $USER \
        && mkdir -p /etc/sudoers.d \
        && echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER \
        && chmod 0440 /etc/sudoers.d/$USER

EXPOSE 3000

ENTRYPOINT ["/worklog"]