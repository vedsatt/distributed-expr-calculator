package orchestrator

// import (
// 	"log"
// 	"net"

// 	pb "github.com/vedsatt/calc_prl/api/gen/go"
// 	"github.com/vedsatt/calc_prl/internal/models"
// 	"google.golang.org/grpc"
// )

// type Server struct {
// 	pb.OrchestratorServer
// }

// func NewServer() *Server {
// 	return &Server{}
// }

// func (s *Server) Calculate(stream pb.Orchestrator_CalculateServer) {
// 	go func() {
// 		for {
// 			select {
// 			case task := <-tasksCh:
// 				grpcTask := &pb.TaskRequest{
// 					ID:       int32(task.ID),
// 					Arg1:     task.Left.Value,
// 					Arg2:     task.Right.Value,
// 					Operator: task.AstType,
// 				}
// 				if err := stream.Send(grpcTask); err != nil {
// 					log.Printf("Failed to send task: %v", err)
// 				}
// 			}
// 		}
// 	}()

// 	go func() {
// 		res, err := stream.Recv()
// 		if err != nil {
// 			log.Printf("Failed to get receive from agent: %v", err)
// 			result := models.Result{
// 				ID:     int(res.ID),
// 				Result: float64(res.Result),
// 				Error:  res.Error,
// 			}
// 			resultsCh <- result
// 		}
// 	}()
// }

// func runGRPC() {
// 	addr := "localhost:5000"
// 	lis, err := net.Listen("tcp", addr)
// 	if err != nil {
// 		log.Fatalf("error starting tcp server: %v", err)
// 	}

// 	log.Printf("tcp server started on port: %v", port)
// 	grpcServer := grpc.NewServer()
// 	orchServer := NewServer()
// 	pb.RegisterOrchestratorServer(grpcServer, orchServer.OrchestratorServer)

// 	if err := grpcServer.Serve(lis); err != nil {
// 		log.Fatalf("error serving grpc: %v", err)
// 	}
// }
