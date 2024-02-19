// docker-compose up -d && sleep 3 && docker exec -it apg_team00 psql -U pitermar -d anomaly
// select * from anomalies;

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	device "server/device"
	"server/report"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address = "localhost:50051"
	caFile  = "ex00/minica/minica.pem"
)

type Storage struct {
	a [100]float64
}

var pool = sync.Pool{
	New: func() interface{} { return new(Storage) },
}

func main() {
	var n = flag.Float64("k", 2, "Anomaly coefficient")
	flag.Parse()
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	stream, err := c.Connect(ctx, &device.ConnectRequest{})
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	var Data *Storage
	Data = pool.Get().(*Storage)
	var sum = 0.0
	var i = 0
	for i < 100 {
		data, err := stream.Recv()
		if err == io.EOF {
			os.Exit(1)
		}
		if err != nil {
			log.Fatalf("Failed to receive: %v", err)
		}
		((*Data).a)[i] = data.Frequency
		sum += data.Frequency
		i++
		if i%10 == 0 {
			log.Printf("10 lines processed \n")
		}
	}
	var mean = sum / 100
	var devSum = 0.0
	i--
	for i >= 0 {
		devSum += (((*Data).a)[i] - mean) * (((*Data).a)[i] - mean)
		i--
	}
	pool.Put(Data)
	devSum /= 100
	fmt.Printf("\nMean = %f\n", mean)
	fmt.Printf("STD = %f\n", devSum)
	var Max = mean + (*n)*devSum
	var Min = mean - (*n)*devSum
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			os.Exit(1)
		}
		if err != nil {
			log.Fatalf("Failed to receive: %v", err)
		}
		if data.Frequency > Max || data.Frequency < Min {
			log.Printf("Anomaly: %s %f %v", data.SessionId, data.Frequency, data.UtcTimestamp)

			err = report.Report(data.SessionId, data.Frequency, data.UtcTimestamp)
			if err != nil {
				log.Fatalf("Failed to saving: %v", err)
			}
		}

	}

}
