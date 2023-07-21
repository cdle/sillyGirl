package core

import (
	"log"
	"net"

	"github.com/cdle/sillyGirl/proto3/srpc"
	"google.golang.org/grpc"
)

// protoc --go_out=. -I. --go-grpc_out=.  bucket.proto
func init() {
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		srpc.RegisterSillyGirlServiceServer(s, &SillyGirlService{})
		log.Printf("grpc server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}
