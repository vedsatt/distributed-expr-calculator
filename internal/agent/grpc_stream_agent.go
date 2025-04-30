package agent

import (
	pb "github.com/vedsatt/calc_prl/api/gen/go"
)

type Server struct {
	pb.CalculatorServer
}

func (s *Server) Calculate(req *pb.ExpressionRequest, stream pb.Orchestrator_CalculateServer) {

}
