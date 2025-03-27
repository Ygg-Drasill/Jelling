package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	pb "github.com/Ygg-Drasill/Jelling/service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strings"
)

type T struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Sha         string `json:"sha"`
	Size        int    `json:"size"`
	Url         string `json:"url"`
	HtmlUrl     string `json:"html_url"`
	GitUrl      string `json:"git_url"`
	DownloadUrl string `json:"download_url"`
	Type        string `json:"type"`
	Content     string `json:"content"`
	Encoding    string `json:"encoding"`
	Links       struct {
		Self string `json:"self"`
		Git  string `json:"git"`
		Html string `json:"html"`
	} `json:"_links"`
}

var (
	network = "localhost"
	address = ":50051"
)

func main() {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s%s", network, address),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			log.Fatalf("Could not close: %v", err)
		}
	}(conn)

	client := pb.NewFileClient(conn)

	req := &pb.FileRequest{Owner: "Ygg-Drasill", Repo: "Jelling", Path: "README.md"}
	resp, err := client.FetchFile(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not fetch file: %v", err)
	}

	j := T{}
	err = json.Unmarshal(resp.Content, &j)
	if err != nil {
		log.Fatalf("Could not unmarshal: %v", err)
	}
	j.Content = strings.ReplaceAll(j.Content, "\n", "")

	d, err := base64.StdEncoding.DecodeString(j.Content)
	if err != nil {
		log.Fatalf("Could not decode: %v", err)
	}
	fmt.Printf("%s\n", d)
}
