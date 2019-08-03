package main

import (
	"context"
	"log"
	"net"

	pb "gogithub/protos"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type helloworldserver struct{}
type githubserver struct{}

func (s *githubserver) FetchByUsername(ctx context.Context, in *pb.GithubRequest) (*pb.GithubResponse, error) {
	log.Printf("[GithubServer] Received Username: %v", in.Username)
	return &pb.GithubResponse{Username: in.Username}, nil
}

func (s *helloworldserver) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("[HelloWorld] Received: %v", in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHelloWorldServer(s, &helloworldserver{})
	pb.RegisterGithubServiceServer(s, &githubserver{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
