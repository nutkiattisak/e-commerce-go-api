# Build stage
FROM golang:1.25.1-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Bangkok

RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

COPY --from=builder /app/main .

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["./main"]
