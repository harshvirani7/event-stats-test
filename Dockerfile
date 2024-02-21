# Dockerfile for building and running the Go service

FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o event-stats-test .

EXPOSE 8080

CMD ["./event-stats-test"]
