package main

import (
	pb "blob-service/gen/blob/v1"
	"blob-service/internal/config"
	"blob-service/internal/server"
	"blob-service/internal/storage"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	conf := config.Load()

	db, err := storage.NewDB(conf.DBHost, conf.DBPort, conf.DBUser, conf.DBPassword, conf.DBName)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	repo := storage.NewBlobRepository(db)
	blobServer := server.NewBlobServer(repo)

	listener, err := net.Listen("tcp", ":"+conf.ServerPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBlobServiceServer(grpcServer, blobServer)

	log.Printf("gRPC server is listening on %s", conf.ServerPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
