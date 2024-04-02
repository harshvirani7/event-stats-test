FROM golang:1.22.0-alpine AS builder

ENV HOME /usr/src/app
WORKDIR $HOME

# Copy the Go application source code into the container
COPY . .

RUN go mod download

# Build the Go application
RUN cd cmd/search && \
    go build -o search .

# Start a new stage from scratch
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the previous stage
COPY --from=builder /usr/src/app/cmd/search/search .

# Copy the production.json file into the /config/search directory
COPY ./config/search/production.json /config/search/production.json

EXPOSE 8080

CMD ["./search"]
