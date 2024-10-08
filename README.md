# Online Canvas Games

## Overview

**Online Canvas Games** is a fullstack web application to provide a way to play some games online.

## Main Features

- **Multiplayer Online Games**: Play with friends or meet new ones while enjoying a collection of fun and interactive games.
- **Real-Time Interactions**: Experience seamless gameplay with real-time updates and interactions among players.
- **User Authentication**: Secure login system to protect user accounts and personal information.

## Technologies Used

- **Backend**: [`Go`](https://golang.org/)
- **Routing**: [`go-chi/chi`](https://github.com/go-chi/chi)
- **Authentication**: [`golang-jwt/jwt`](https://github.com/golang-jwt/jwt).
- **WebSocket Communication**: [`gorilla/websocket`](https://github.com/gorilla/websocket).
- **Database**: [`PostgreSQL`](https://www.postgresql.org/), accessed via [`lib/pq`](https://github.com/lib/pq) and administrated with [`adminer`](https://www.adminer.org/).
- **Cryptography**: [`golang.org/x/crypto`](https://golang.org/x/crypto).
- **Frontend**: [`TypeScript`](https://www.typescriptlang.org/).
- **TypeScript Recursive Compilation**: [`tscompiler`](https://github.com/1001bit/tscompiler)
- **Containerization**: [`Docker`](https://www.docker.com/).
- **Inter-Microservice Communication**: [`gRPC`](https://grpc.io/) 
- **Templating**: [`a-h/templ`](https://github.com/a-h/templ/)

## Getting Started

1. To get started with **Online Canvas Games**, ensure you have Docker installed on your machine.
2. Clone the repository and navigate to the project directory.
3. Make sure environmental variables are set (see .env)
4. If want to edit typescript files, make sure `tsc` tool is installed
5. If want to edit .templ files, make sure `templ` utility is installed (`go install github.com/a-h/templ/cmd/templ@latest`)

## Build and start

- Build everything and start: `make all`
- Build services: `make build`
- Start compose: `make up`
- Stop compose: `make down`
- Compile typescript files: `make compile-ts`
- Compile .templ files: `make compile-templ`
- Clean up docker resources: `make clean`

## Microservices

- **gateway**: 
  - Manages API endpoints and web page serving
  - Handles user authorization
  - Routes requests to appropriate microservices

- **games**: 
  - Implements game logic
  - Handles WebSocket connections for real-time gameplay
  - Manages game data and rooms state

- **user**: 
  - Manages user authentication
  - Stores and retrieves user profile data

- **storage**: 
  - Manages static assets and file storage
  - Handles image retrievals
  - Compiles TypeScript to JavaScript

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.
