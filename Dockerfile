FROM golang:1.21.4-alpine3.18 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /mailer

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o /app/main cmd/main.go
COPY /config/config.yml /app/config/config.yml


FROM alpine:3

WORKDIR /app
COPY --from=builder /app/main /app/main
COPY --from=builder /app/config/config.yml /app/config/config.yml

CMD ["./main"]