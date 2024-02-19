// brew install protobuf
// brew install protoc-gen-go

// go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
// export PATH="$PATH:$(go env GOPATH)/bin"
// go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.0.0

// go get google.golang.org/protobuf/reflect/protoreflect
// go get google.golang.org/grpc
// go get github.com/gofrs/uuid
// protoc --go_out=. --go-grpc_out=. device/device.proto
// go mod tidy

package main

import (
	"crypto/tls"
	"log"
	"math/rand"
	"net"
	"time"

	device "server/device"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type server struct {
	device.UnimplementedDeviceServiceServer
}

func (s *server) Connect(_ *device.ConnectRequest, stream device.DeviceService_ConnectServer) error {
	sessionID, _ := uuid.NewV4()

	mean := rand.Float64()*20 - 10           // рандом float [-10, 10]
	stdDev := rand.Float64()*(1.5-0.3) + 0.3 // рандом float [0.3, 1.5]

	log.Printf("New connection: session_id: %s, mean: %.2f, stddev: %.2f", sessionID, mean, stdDev)

	for {
		frequency := rand.NormFloat64()*stdDev + mean

		data := &device.DeviceData{
			SessionId:    sessionID.String(),
			Frequency:    frequency,
			UtcTimestamp: time.Now().Unix(),
		}

		if err := stream.Send(data); err != nil {
			return err
		}

		time.Sleep(time.Second)
	}
}

const (
	port     = "localhost:50051"
	certFile = "ex00/minica/localhost/cert.pem"
	keyFile  = "ex00/minica/localhost/key.pem"
)

func main() {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}
	// Создаем объект creds, который содержит ключи сервера
	creds := credentials.NewServerTLSFromCert(&cert)

	s := grpc.NewServer(grpc.Creds(creds))
	device.RegisterDeviceServiceServer(s, &server{})

	log.Printf("Server started on port: %s", port)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	// подключение reflection для тестирования: "evans -r repl -p 50051"
	// device.DeviceService@127.0.0.1:50051> show service
	reflection.Register(s) // device.DeviceService@127.0.0.1:50051> call Connect

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
