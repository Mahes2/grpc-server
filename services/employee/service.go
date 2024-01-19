package employee

import (
	"context"

	"github.com/grpc-server/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.EmployeeServer
}

func (s *Server) GetById(ctx context.Context, req *pb.GetByIdRequest) (*pb.GetByIdResponse, error) {
	switch req.Id {
	case 1000:
		return nil, status.Error(codes.InvalidArgument, "id can't be 1000")
	case 0:
		panic("id can't be 0")
	default:
		return &pb.GetByIdResponse{
			Id:   req.Id,
			Name: "John Doe",
		}, nil
	}
}
