package stream

import (
	"crl/stream/streampb"
	"fmt"
	"io"
	"net"
	"time"

	log "github.com/inconshreveable/log15"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServerHandler struct {
}

func (s *ServerHandler) SubscribeToEvent(stream streampb.StreamService_SubscribeToEventServer) error {
	// infinite loop
	go func() {
		for {
			time.Sleep(5 * time.Second)
			err := stream.Send(&streampb.StreamResponse{
				Status: "OK",
				Payload: "hello",
			})
			if err != nil {
				log.Error("Failed to send", "err", err)
			}
		}
	}()

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error("Failed to receive", "err", err)
		}
		log.Info("Got", "code", in.Type, "payload", in.Payload)
	}


	return nil
}

func Start() {
	h := ServerHandler{}

	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Error("failed to listen", "ERROR", err)
		panic(fmt.Sprintf("failed to listen: %v", err))
	}

	// Create and attach a new the gprc service to a new ServerHandler
	grpcServer := grpc.NewServer()
	streampb.RegisterStreamServiceServer(grpcServer, &h)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Error("failed to start grpc auth service", "ERROR", err)
		panic(fmt.Sprintf("failed to serve: %s", err))
	}

}
