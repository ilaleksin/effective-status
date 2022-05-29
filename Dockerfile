FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"

FROM alpine

COPY --from=builder /app/effective-status /usr/bin

WORKDIR /app

ENTRYPOINT ["effective-status"]
