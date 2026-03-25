package main

import (
	"context"
	"log"
	"net/http"

	pb "grpc-crud/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	err := pb.RegisterUserServiceHandlerFromEndpoint(
		ctx, 
		mux,
		"localhost:50051",
		opts,
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("REST gateway running on :8080")
	http.ListenAndServe(":8080", mux)
}