
FROM golang:1.22-alpine AS builder


RUN apk add --no-cache git

WORKDIR /workspace


COPY app/go.mod app/go.sum ./
RUN go mod download


COPY app/ .


RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /gitcrawler ./main/


FROM alpine:3.20

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY --from=builder /gitcrawler .

VOLUME ["/app"]

EXPOSE 8080

CMD ["./gitcrawler"]
