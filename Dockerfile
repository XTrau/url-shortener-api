FROM golang:1.26-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o main cmd/main.go

FROM alpine:latest

COPY --from=builder /app .
EXPOSE 8080
CMD ["./main"]
