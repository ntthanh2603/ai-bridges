# Stage 1: build
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go

# Stage 2: Final Image
FROM alpine:3.22.2

WORKDIR /root/

COPY --from=builder /app/main .

ENTRYPOINT ["./main"] 
