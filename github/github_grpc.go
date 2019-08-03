package github

import (
	"context"
	"encoding/json"
	pb "gogithub/protos"
	"log"
)

// GrpcServer is github grpc server
type GrpcServer struct{}

// FetchByUsername = implement from proto
func (s *GrpcServer) FetchByUsername(ctx context.Context, in *pb.GithubRequest) (*pb.GithubResponse, error) {
	log.Printf("[GrpcServer] Received Username: %v", in.Username)
	data, err := FetchAllRepos(in.Username)
	if err != nil {
		log.Fatalf("failed to fetch github: %v", err)
	}

	langmap := make(map[string]int32)
	b, _ := json.Marshal(data.LanguageMap)
	json.Unmarshal(b, &langmap)

	return &pb.GithubResponse{
		Username:        in.Username,
		Starcount:       int32(data.StarCount),
		Repocount:       int32(data.RepoCount),
		Forkcount:       int32(data.ForkCount),
		Watchercount:    int32(data.WatcherCount),
		Subscribercount: int32(data.SubscriberCount),
		Langmap:         langmap,
	}, nil
}
