FROM golang:1.23-alpine3.21 as builder

WORKDIR /obscura

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./bin/obscura ./

FROM alpine:3.21 as app

WORKDIR /app

COPY --from=builder /obscura/bin/obscura .

ENTRYPOINT ["/app/obscura"]
