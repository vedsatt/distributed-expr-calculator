package orchestrator

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/vedsatt/calc_prl/api/gen/go"
	"github.com/vedsatt/calc_prl/internal/models"
	"google.golang.org/grpc"
)

const (
	tcp         = "tcp"
	addr string = ":5000"
)

type Server struct {
	pb.UnimplementedOrchestratorServer
	mu sync.Mutex
}

func NewServer() *Server {
	return &Server{mu: sync.Mutex{}}
}

func (s *Server) Calculate(stream pb.Orchestrator_CalculateServer) error {
	log.Println("agent connected to gRPC server")
	defer log.Println("agent disconnected")
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	done := make(chan struct{})
	defer close(done)

	go func() {
		defer cancel()
		for {
			select {
			case task := <-tasksCh:
				s.mu.Lock()
				err := stream.Send(&pb.TaskRequest{
					Id:       int32(task.ID),
					Arg1:     task.Left.Value,
					Arg2:     task.Right.Value,
					Operator: task.Value,
				})
				s.mu.Unlock()

				if err != nil {
					log.Printf("Failed to send task: %v", err)
					return
				}
			case <-ctx.Done():
				return
			case <-done:
				return
			}
		}
	}()

	go func() {
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				res, err := stream.Recv()
				if err != nil {
					log.Printf("Receive error: %v", err)
					return
				}
				resultsCh <- models.Result{
					ID:     int(res.Id),
					Result: float64(res.Result),
					Error:  res.Error,
				}
			}
		}
	}()

	<-ctx.Done()
	return nil
}

func runGRPC() {
	log.Println("Starting tcp server...")
	lis, err := net.Listen(tcp, addr)
	if err != nil {
		log.Fatalf("error starting tcp server: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrchestratorServer(grpcServer, NewServer())

	log.Printf("tcp server started at: %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("error serving grpc: %v", err)
	}
}
