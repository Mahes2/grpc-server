package employee

import (
	"context"

	"github.com/grpc-server/pb"
)

type Server struct {
	pb.EmployeeServer
}

func (s *Server) GetById(ctx context.Context, req *pb.GetByIdRequest) (*pb.GetByIdResponse, error) {
	return &pb.GetByIdResponse{
		Id:   req.Id,
		Name: "John Doe",
	}, nil
}
