FROM golang:1.22.5 AS builder

WORKDIR /app

# dependencies
COPY go.* .
RUN go mod download
# other files
COPY . .

# compiling ts -> js
RUN apt update && apt install node-typescript -y
RUN go run cmd/tscompiler/main.go

# compiling go -> ./main
RUN go build -o main cmd/server/main.go

# Final stage
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/main ./main
COPY --from=builder /app/web ./web