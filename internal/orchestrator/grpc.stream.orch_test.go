package orchestrator

import (
	"context"
	"net"
	"testing"
	"time"

	pb "github.com/vedsatt/calc_prl/api/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestServerConnection(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterOrchestratorServer(srv, NewServer())
	go srv.Serve(lis)
	defer srv.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrchestratorClient(conn)

	stream, err := client.Calculate(ctx)
	if err != nil {
		t.Fatalf("Calculate failed: %v", err)
	}
	defer stream.CloseSend()
}

func TestOrchestratorGRPC(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterOrchestratorServer(srv, NewServer())
	go srv.Serve(lis)
	defer srv.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrchestratorClient(conn)

	t.Run("successful task processing", func(t *testing.T) {
		stream, err := client.Calculate(ctx)
		if err != nil {
			t.Fatalf("Calculate failed: %v", err)
		}

		if err := stream.Send(&pb.AgentResponse{
			Id:     1,
			Result: 42.0,
		}); err != nil {
			t.Fatalf("Send failed: %v", err)
		}
	})
}
