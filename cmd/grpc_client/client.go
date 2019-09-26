package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	pb "gogithub/protos"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	cgithub := pb.NewGithubServiceClient(conn)

	// Contact the server and print out its response.
	var username string
	if len(os.Args) > 1 {
		username = os.Args[1]
	} else {
		log.Fatalf("Please provide username as argument")
	}
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()
	ctx := context.Background()

	resGithub, err := cgithub.FetchByUsername(ctx, &pb.GithubRequest{Username: username})
	if err != nil {
		log.Fatalf("could not fetch: %v", err)
	}
	log.Printf("[GithubGrpcClient]: %s (%d stars, %d repos, %d forks)\n",
		resGithub.Username,
		resGithub.Starcount,
		resGithub.Repocount,
		resGithub.Forkcount)
	b, _ := json.MarshalIndent(resGithub.Langmap, "", "  ")
	log.Printf("LangMap: %s", string(b))
}
