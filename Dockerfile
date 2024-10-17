FROM golang:1.23-alpine AS builder

COPY . /github.com/solumD/auth
WORKDIR /github.com/solumD/auth

RUN go mod download
RUN go build -o ./bin/auth_server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/solumD/auth/bin/auth_server .

ADD .env .

CMD ["./auth_server"]