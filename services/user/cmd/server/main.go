package main

import (
	"fmt"
	"log"
	"net"

	"github.com/1001bit/onlinecanvasgames/services/user/internal/database"
	"github.com/1001bit/onlinecanvasgames/services/user/internal/server"
	"github.com/1001bit/onlinecanvasgames/services/user/pkg/env"
	"github.com/1001bit/onlinecanvasgames/services/user/pkg/userpb"
	"google.golang.org/grpc"
)

func main() {
	// start database
	err := database.Start()
	if err != nil {
		log.Fatal("err starting database:", err)
	}
	defer database.DB.Close()

	// start listener
	addr := fmt.Sprintf(":%s", env.GetEnvVal("PORT"))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	// create server
	s := grpc.NewServer()
	userpb.RegisterUserServiceServer(s, server.NewUserServer())

	// serve listener
	log.Println("listening and", listener.Addr())
	log.Fatal(s.Serve(listener))
}
