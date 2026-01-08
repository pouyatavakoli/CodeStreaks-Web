# ---------- Builder stage ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk update && apk add --no-cache git ca-certificates

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build main binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api cmd/api/main.go

# Build add_user script binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o add_user scripts/add_user.go

# ---------- Runtime stage ----------
FROM alpine:3.19

WORKDIR /app

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk update && apk add --no-cache ca-certificates tzdata

# Copy compiled binaries
COPY --from=builder /app/api .
COPY --from=builder /app/add_user .

# Copy frontend
COPY --from=builder /app/frontend ./frontend

EXPOSE 8080
CMD ["./api"]

