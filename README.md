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
- **Database**: PostgreSQL, accessed via [`lib/pq`](https://github.com/lib/pq).
- **Environment Variables**: [`joho/godotenv`](https://github.com/joho/godotenv).
- **Cryptography**: [`golang.org/x/crypto`](https://golang.org/x/crypto).

## Getting Started

To get started with **Online Canvas Games**, ensure you have Go installed on your system. Then, clone the repository and navigate to the project directory.


Before running the server, make sure to set up your PostgreSQL database and environment variables. You can use `.env.example` as a template for creating your own `.env` file.

To compile TypeScript files for the client-side (Linux machine):
` ./scripts/tscompiler.sh `


To run the server (Linux machine):
`./scripts/run.sh`


## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.