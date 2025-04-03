package main

import (
	"context"
	"fmt"
	pb "github.com/Ygg-Drasill/Jelling/service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"net/http"
)

type Server struct {
	pb.UnimplementedFileServer
}

var (
	network = "tcp"
	address = ":50051"
)

func (s *Server) FetchFile(ctx context.Context, req *pb.FileRequest) (*pb.FileResponse, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", req.Owner, req.Repo, req.Path)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("failed to close body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch file: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("request received: %s", url)
	return &pb.FileResponse{Content: body}, nil
}

func main() {
	listen, err := net.Listen(network, address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterFileServer(s, &Server{})
	reflection.Register(s)
	log.Printf("Listening on %s%s", network, address)
	if err = s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
