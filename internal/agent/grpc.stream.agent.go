package agent

import (
	"context"
	"log"
	"time"

	pb "github.com/vedsatt/calc_prl/api/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Connect() {
	for {
		conn, err := grpc.NewClient(
			"orchestrator:5000",
			grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			log.Printf("error connecting to the server: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		client := pb.NewOrchestratorClient(conn)
		err = handleStream(client)
		conn.Close()

		if err != nil {
			log.Printf("stream error: %v", err)
		}
		time.Sleep(1 * time.Second)
	}

}

func handleStream(client pb.OrchestratorClient) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.Calculate(ctx)
	if err != nil {
		return err
	}

	done := make(chan struct{})
	defer close(done)

	go func() {
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				return
			case <-done:
				return
			default:
				task, err := stream.Recv()
				if err != nil {
					log.Printf("Receive error: %v", err)
					return
				}

				tasksCh <- &Task{
					ID:   int(task.Id),
					Arg1: task.Arg1,
					Arg2: task.Arg2,
					Type: task.Operator,
				}
			}
		}
	}()

	go func() {
		defer cancel()
		for {
			select {
			case result := <-resultsCh:
				err := stream.Send(&pb.AgentResponse{
					Id:     int32(result.ID),
					Result: float32(result.Result),
					Error:  result.Error,
				})
				if err != nil {
					log.Printf("Send error: %v", err)
					return
				}
			case <-ctx.Done():
				return
			case <-done:
				return
			}
		}
	}()

	<-ctx.Done()
	return ctx.Err()
}
