FROM golang:1.24-alpine

WORKDIR /app

COPY . .

RUN go mod download && \
    go mod verify

RUN GOOS=linux GOARCH=amd64 \
    go build -o /app/bin/subscriptions /app/cmd/subscriptions/main.go

EXPOSE 8080

CMD ["/app/bin/subscriptions"]
