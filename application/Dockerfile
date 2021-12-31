FROM golang:1.16-alpine as builder

RUN apk add binutils

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /webhooks

FROM alpine:latest

WORKDIR /

COPY --from=builder /webhooks /webhooks

EXPOSE 9000

CMD ["/webhooks"]