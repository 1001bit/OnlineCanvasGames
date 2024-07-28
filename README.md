# Online Canvas Games

## Overview

**Online Canvas Games** is a web application designed to provide a platform for users to enjoy multiplayer online games together. It features real-time interactions, user authentication, and a variety of engaging games that players can join and play with others across the globe.

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
- **Environment Variables**: [`joho/godotenv`](https://github.com/joho/godotenv).
- **Cryptography**: [`golang.org/x/crypto`](https://golang.org/x/crypto).
- **Frontend**: [`TypeScript`](https://www.typescriptlang.org/).
- **Containerization**: [`Docker`](https://www.docker.com/).

## Getting Started

1. To get started with **Online Canvas Games**, ensure you have Docker installed on your machine.
2. Clone the repository and navigate to the project directory.
3. Create your own `.env` in the root directory based on the `example.env`

## Running and building the Application

1. To build the entire application, use `docker compose build`
2. To run, use `docker compose up` 

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.