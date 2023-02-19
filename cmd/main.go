package main

import (
	"fmt"
	"log"
	"net"

	"github.com/dedihartono801/auth-svc/pkg/config"
	"github.com/dedihartono801/auth-svc/pkg/db"
	"github.com/dedihartono801/auth-svc/pkg/services"
	"github.com/dedihartono801/auth-svc/pkg/utils"
	pb "github.com/dedihartono801/protobuf/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := db.Init(c.DBUrl)

	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "go-grpc-auth-svc",
		ExpirationHours: 24 * 365,
	}

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Auth Svc on", c.Port)

	s := services.Server{
		H:   h,
		Jwt: jwt,
	}

	opts := []grpc.ServerOption{}
	tls := true

	if tls {
		certFile := "ssl/auth-svc/server.crt"
		kefFile := "ssl/auth-svc/server.pem"

		creds, err := credentials.NewServerTLSFromFile(certFile, kefFile)

		if err != nil {
			log.Fatalf("Failed loading certificates: %v\n", err)
		}

		opts = append(opts, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
