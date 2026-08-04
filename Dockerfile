
FROM golang:1.22-alpine AS builder


RUN apk add --no-cache git

WORKDIR /workspace


COPY app/ .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /gitcrawler ./main/


FROM alpine:3.20

RUN apk add --no-cache git ca-certificates

COPY --from=builder /gitcrawler /usr/local/bin/gitcrawler

RUN mkdir -p /app
WORKDIR /app

VOLUME ["/app"]

EXPOSE 8080

CMD ["gitcrawler"]
