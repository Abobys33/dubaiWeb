FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -mod=vendor -o /dubai-web ./cmd/dubai-web

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /dubai-web /app/dubai-web
COPY ./web/static ./web/static

RUN chmod +x /app/dubai-web

EXPOSE 8080

CMD ["/app/dubai-web"]