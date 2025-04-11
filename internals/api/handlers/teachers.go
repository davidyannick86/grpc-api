package handlers

import (
	"context"

	"github.com/davidyannick86/grpc-api-mongodb/internals/repositories/mongodb"
	pb "github.com/davidyannick86/grpc-api-mongodb/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AddTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {

	for _, teacher := range req.GetTeachers() {
		if teacher.Id != "" {
			return nil, status.Error(codes.InvalidArgument, "Id must be empty")
		}
	}

	addedTeachers, err := mongodb.AddTeacherToDb(ctx, req.GetTeachers())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Teachers{Teachers: addedTeachers}, nil
}
