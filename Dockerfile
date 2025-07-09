FROM golang:1.24.3-slim

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd

CMD ["./main"]
