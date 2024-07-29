FROM golang:1.22.5 AS builder

WORKDIR /app

# update and install tsc
RUN apt update && apt install -y node-typescript

# golang dependencies
COPY go.* .
RUN go mod download

# other files
COPY . .

# compiling go -> ./main
RUN go build -o main cmd/server/main.go

# compiling ts -> js
RUN go run cmd/tscompiler/main.go

# Final stage
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/main ./main
COPY --from=builder /app/storage ./storage