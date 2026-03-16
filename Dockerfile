# ── Stage 1: Build ────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

# git is needed by `go mod download` for VCS-backed modules
RUN apk add --no-cache git

WORKDIR /workspace

# Copy module manifests first so that dependency layer is cached
COPY app/go.mod app/go.sum ./
RUN go mod download

# Copy the rest of the source
COPY app/ .

# Build a static binary – no CGO needed for this service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /gitcrawler ./main/

# ── Stage 2: Runtime ──────────────────────────────────────────────────────
FROM alpine:3.20

# git is required at runtime because the service runs `git clone`
# ca-certificates lets the binary reach HTTPS endpoints (OpenRouter, GitHub)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY --from=builder /gitcrawler .

# Repos are cloned to a sub-directory of the working directory
# (CloneService calls os.Getwd() then os.MkdirTemp).
# Mount a named volume here for persistence and to avoid container bloat.
VOLUME ["/app"]

EXPOSE 8080

CMD ["./gitcrawler"]
