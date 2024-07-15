# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS builder


ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go  build -v github.com/oyamo/forumz-auth-server/cmd/server

FROM gcr.io/distroless/base
COPY --from=builder /app/server /

EXPOSE 3000
CMD ["/server"]
