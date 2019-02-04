package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/innovate-technologies/geo-service/pb"
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
	c := pb.NewGeoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetGeoInfo(ctx, &pb.GeoInfoRequest{Ip: os.Args[1]})
	if err != nil {
		log.Fatalf("could not get geo info: %v", err)
	}
	spew.Dump(r)
}
