package aitools

import (
	"context"
	"github.com/ceerdecy/nautilus-proto-go/core/aitool/pb"
)

func (p *provider) ListAiTools(ctx context.Context, request *pb.ListAiToolsRequest) (*pb.ListAiToolsResponse, error) {
	return &pb.ListAiToolsResponse{}, nil
}
