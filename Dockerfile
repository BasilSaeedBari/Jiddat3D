# Build stage
FROM golang:latest AS builder
WORKDIR /app

# Download Tailwind CLI v3 (v4 requires different config syntax)
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.17/tailwindcss-linux-x64 \
    && chmod +x tailwindcss-linux-x64 \
    && mv tailwindcss-linux-x64 tailwindcss

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build CSS
RUN ./tailwindcss -i ui/static/css/input.css -o ui/static/css/output.css --minify

# Build Go binary
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /jiddat3d ./cmd/jiddat

# Final stage
FROM alpine:3.19
WORKDIR /app

# Install runtime dependencies (cwebp for image processing, ca-certificates for TLS, mailcap for mime types)
RUN apk add --no-cache libwebp-tools ca-certificates tzdata mailcap

COPY --from=builder /jiddat3d /app/jiddat3d
COPY --from=builder /app/ui /app/ui

# PocketBase data volume
VOLUME /app/pb_data

EXPOSE 8080

CMD ["/app/jiddat3d", "serve", "--http=0.0.0.0:8080"]
