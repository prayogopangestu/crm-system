# Backend Dockerfile for Railway
FROM golang:1.25-alpine AS builder

WORKDIR /app/backend

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/backend/api .
COPY --from=builder /app/backend/configs ./configs

EXPOSE 8080

CMD ["./api"]
