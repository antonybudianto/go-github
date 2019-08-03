package main

import (
	"context"
	"log"
	"net"

	"gogithub/github"
	pb "gogithub/protos"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type helloworldserver struct{}

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
	pb.RegisterGithubServiceServer(s, &github.GrpcServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
