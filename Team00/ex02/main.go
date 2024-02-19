package main

import (
	"context"
	"io"
	"log"
	"time"

	device "ex02_report/device"
	"ex02_report/report"

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
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	c := device.NewDeviceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	stream, err := c.Connect(ctx, &device.ConnectRequest{})
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	for i := 0; i < 10; i++ {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to receive: %v", err)
		}

		log.Printf("Received: %v", data)
		err = report.Report(data.SessionId, data.Frequency, data.UtcTimestamp)
		if err != nil {
			log.Fatalf("Failed to saving: %v", err)
		}
	}

	log.Println("***************************************")
	log.Println()
	log.Println("------ ALL RESULTS FROM POSTGRES ------")
	log.Println()
	log.Println("***************************************")
	err = report.ReadReport()
	if err != nil {
		log.Fatalf("Failed to reading report: %v", err)
	}
}
