FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api cmd/api/main.go

FROM alpine:3.22.2

RUN apk --no-cache add ca-certificates wget

WORKDIR /root/

COPY --from=builder /app/api .
COPY --from=builder /app/data ./data
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./api"]
