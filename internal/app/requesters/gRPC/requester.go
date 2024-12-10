package gRPC

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"io"
	"os"
	"yandex_GophKeeper_client/internal/app/requesters/gRPC/proto"
)

const mdJWTKey = "jwt"

// GRPCGophKeeperRequester can send request to a gRPC server.
type GRPCGophKeeperRequester struct {
	C   proto.GophKeeperServiceClient
	JWT string
	// MaxBinDataChunkSize - value in bytes.
	MaxBinDataChunkSize int
	Logger              *zap.SugaredLogger
}

func NewGRPCRequester(client proto.GophKeeperServiceClient, jwt string, maxBinDataChunkSize int, logger *zap.SugaredLogger) *GRPCGophKeeperRequester {
	return &GRPCGophKeeperRequester{
		C:                   client,
		JWT:                 jwt,
		MaxBinDataChunkSize: maxBinDataChunkSize,
		Logger:              logger,
	}
}

func (g *GRPCGophKeeperRequester) SendBinFile(path string, dataName string) error {
	if len(path) == 0 {
		return fmt.Errorf("path is empty")
	}
	if len(dataName) == 0 {
		return fmt.Errorf("dataName is empty")
	}
	// open file
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	//put jwt into ctx
	ctx := metadata.AppendToOutgoingContext(context.Background(), mdJWTKey, g.JWT)

	stream, err := g.C.SaveBinData(ctx)
	if err != nil {
		return fmt.Errorf("failed to initiate SaveBinData stream: %w", err)
	}

	//starting send data
	buffer := make([]byte, g.MaxBinDataChunkSize)
	for {
		//read chunk
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		//prepare chunk
		reqChunk := &proto.SaveBinDataRequest{
			DataName: dataName,
			Chunk:    buffer[:n],
		}

		// send chunk
		if err := stream.Send(reqChunk); err != nil {
			return fmt.Errorf("failed to send chunk: %w", err)
		}
	}

	//close stream
	if _, errCloseSteam := stream.CloseAndRecv(); errCloseSteam != nil && errCloseSteam != io.EOF {
		return fmt.Errorf("failed to close stream: %w", err)
	}

	return nil
}

func (g *GRPCGophKeeperRequester) GetBinFile(fileName, outputPath string) error {
	//put jwt into ctx
	ctx := metadata.AppendToOutgoingContext(context.Background(), mdJWTKey, g.JWT)

	//prepare request
	req := &proto.GetBinDataRequest{
		DataName: fileName,
	}

	//start stream
	stream, err := g.C.GetBinData(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to initiate GetBinData stream: %w", err)
	}

	//create file to write
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	//get data from server
	for {
		//read chunk
		resp, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("failed to receive chunk: %w", err)
		}

		//write chink to a file
		if _, err := file.Write(resp.Chunk); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	return nil
}
