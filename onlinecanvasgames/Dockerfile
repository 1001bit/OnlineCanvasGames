FROM golang:1.22.5

# directory
WORKDIR /onlinecanvasgames
COPY . /onlinecanvasgames/

# installing tsc
RUN apt update
RUN apt install node-typescript -y

# compiling ts -> js
RUN go run cmd/tscompiler/main.go

# compiling go -> ./main
RUN go build -o main cmd/server/main.go