package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/prantoran/gcd-grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

// Declare the Compute handler function. This makes the server type conform to the
// auto-generated pb.GCDServiceServer interface.
func (s *server) Compute(ctx context.Context, r *pb.GCDRequest) (*pb.GCDResponse, error) {
	a, b := r.A, r.B
	for b != 0 {
		a, b = b, a%b
	}
	return &pb.GCDResponse{Result: a}, nil
}

func main() {

	// In the main function, register a server type which will handle requests. Then start the gRPC server.
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGCDServiceServer(s, &server{})
	reflection.Register(s)
	fmt.Printf("Serving grpc calls at port 3000\n")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
