FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/billing-service .

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /out/billing-service /app/billing-service

EXPOSE 8081

ENTRYPOINT ["/app/billing-service"]
