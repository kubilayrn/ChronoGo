FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN $(go env GOPATH)/bin/swag init -g cmd/server/main.go -o docs

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata wget
WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]

