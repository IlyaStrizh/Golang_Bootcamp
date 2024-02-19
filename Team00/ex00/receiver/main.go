package main

import (
	"context"
	"io"
	"log"
	"time"

	device "server/device"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address = "localhost:50051"
	caFile  = "ex00/minica/minica.pem"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile(caFile, "")
	if err != nil {
		log.Fatalf("Failed to load credentials: %v", err)
	}

	conn, err := grpc.DialContext(context.Background(), address, grpc.WithTransportCredentials(creds), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := device.NewDeviceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	stream, err := c.Connect(ctx, &device.ConnectRequest{})
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to receive: %v", err)
		}

		log.Printf("Received: %v", data)
	}
}
