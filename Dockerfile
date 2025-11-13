FROM golang:1.25.1 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bot ./cmd/main.go

FROM alpine:3.22.2

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder ./app/bot .
COPY .env.production .env

RUN adduser -D botuser
USER botuser

EXPOSE ${HTTP_PORT}

ENTRYPOINT ["/app/bot"]
