FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY go.su[m] ./
RUN go mod download || true
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/main.go

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/server .
COPY .env .
EXPOSE 8080
CMD ["./server"]
