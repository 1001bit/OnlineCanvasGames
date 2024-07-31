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
- **Containerization**: [`Docker`](https://www.docker.com/).

## Getting Started

1. To get started with **Online Canvas Games**, ensure you have Docker installed on your machine.
2. Clone the repository and navigate to the project directory.
3. Create your own `.env` in the root directory based on the `.env.example` (may keep the values)

## Build and start

1. To build the entire application, run `./scripts/build.sh`
2. To compile typescript into javascript files, run `./scripts/compilets.sh`
3. To run the application, run `./scripts/up.sh`

## Microservices

- **gateway**: 
  - Manages API endpoints and web page serving
  - Handles user authorization
  - Routes requests to appropriate microservices

- **games**: 
  - Implements game logic
  - Handles WebSocket connections for real-time gameplay
  - Manages game data and state

- **user**: 
  - Manages user authentication
  - Stores and retrieves user profile data

- **storage**: 
  - Manages static assets and file storage
  - Handles image retrievals
  - Compiles TypeScript to JavaScript

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.