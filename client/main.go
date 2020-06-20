package main

import (
	"context"
	"log"
	"time"

	"github.com/plaurent-dev/nttbm/pkg/proto/site"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Client running ...")

	conn, err := grpc.Dial(":50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := site.NewSiteServiceClient(conn)

	request := &site.SiteRequest{Url: "https://accounts.google.com"}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.Site(ctx, request)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Site:", response.GetSite())
	log.Printf("Elapsed Time: %d ms", response.GetElapsedtime())
	log.Println("Access:", response.GetAccess())
}
