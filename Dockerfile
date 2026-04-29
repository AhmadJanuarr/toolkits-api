# Stage 1: Build Stage
FROM golang:1.25-alpine AS builder

# Install Compiler C yang dibutuhkan untuk library chai2010/webp
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy dependency
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build aplikasi

# CGO_ENABLED=1 wajib untuk library chai2010/webp
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o main .

# Stage 2: Runtime Stage (Image Akhir)
FROM alpine:latest

WORKDIR /root/

RUN apk update && \
    apk add --no-cache ffmpeg python3 curl nodejs && \
    curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp && \
    chmod a+rx /usr/local/bin/yt-dlp    
  
    
# Copy binary dari builder
COPY --from=builder /app/main .


# Buat folder temp yang dibutuhkan aplikasi
RUN mkdir -p temp/uploads temp/processed temp/compressed temp/resized temp/downloads/youtube temp/downloads/instagram temp/downloads/tiktok

# Expose port
EXPOSE 8080

# Jalankan
CMD ["./main"]
