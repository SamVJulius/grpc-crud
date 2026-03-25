package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "grpc-crud/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
	store map[string]*pb.User
	mu    sync.Mutex
}

func newServer() *server {
	return &server{
		store: make(map[string]*pb.User),
	}
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	user := &pb.User{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
	}

	s.store[id] = user

	return &pb.UserResponse{
		User: user,
	}, nil
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.store[req.Id]
	if !ok {
		return nil, nil
	}

	return &pb.UserResponse{
		User: user,
	}, nil
}

func (s *server) UdpateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.store[req.Id]
	if !ok {
		return nil, nil
	}

	user.Name = req.Name
	user.Email = req.Email

	return &pb.UserResponse{
		User: user,
	}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.store, req.Id)

	return &pb.DeleteUserResponse{
		Message: "user deleted",
	}, nil
}

func (s *server) ListUsers(ctx context.Context, _ *pb.Empty) (*pb.UserList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	users := []*pb.User{}
	for _, u := range s.store {
		users = append(users, u)
	}

	return &pb.UserList{Users: users}, nil
}

func main() {
	lis, _ := net.Listen("tcp", ":50051")

	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, newServer())

	log.Println("server running on port 50051...")

	grpcServer.Serve(lis)
}
