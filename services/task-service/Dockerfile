FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/task-service .

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /out/task-service /app/task-service

EXPOSE 8082

ENTRYPOINT ["/app/task-service"]
