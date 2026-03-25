package main

import (
	"context"
	"log"
	"time"

	pb "grpc-crud/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, 	_ := grpc.Dial("localhost:50051", grpc.WithInsecure())
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, _ := client.CreateUser(ctx, &pb.CreateUserRequest{
		Name: "User",
		Email: "test&email.com",
	})

	log.Println("CreateUser: ", res)

	getRes, _ := client.GetUser(ctx, &pb.GetUserRequest{
		Id: res.User.Id,
	})

	log.Println("Fetched:", getRes.User)
}
