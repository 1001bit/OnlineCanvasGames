package main

import (
	"fmt"
	"log"
	"net"

	"github.com/1001bit/onlinecanvasgames/services/user/internal/database"
	"github.com/1001bit/onlinecanvasgames/services/user/internal/server"
	"github.com/1001bit/onlinecanvasgames/services/user/internal/usermodel"
	"github.com/1001bit/onlinecanvasgames/services/user/pkg/userpb"
	"github.com/1001bit/overenv"
	"google.golang.org/grpc"
)

func main() {
	// start database
	db, err := database.NewFromEnv()
	if err != nil {
		log.Fatal("err starting database:", err)
	}
	defer db.Close()

	// start listener
	addr := fmt.Sprintf(":%s", overenv.Get("PORT"))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	// create server
	s := grpc.NewServer()
	userStore := usermodel.NewUserStore(db)
	userpb.RegisterUserServiceServer(s, server.NewUserServer(userStore))

	// serve listener
	log.Println("listening and", listener.Addr())
	log.Fatal(s.Serve(listener))
}
