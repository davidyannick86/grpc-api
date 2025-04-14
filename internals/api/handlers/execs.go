package handlers

import (
	"context"

	"github.com/davidyannick86/grpc-api-mongodb/internals/models"
	"github.com/davidyannick86/grpc-api-mongodb/internals/repositories/mongodb"
	"github.com/davidyannick86/grpc-api-mongodb/pkg/utils"
	pb "github.com/davidyannick86/grpc-api-mongodb/proto/gen"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AddExecs(ctx context.Context, req *pb.Execs) (*pb.Execs, error) {

	for _, exec := range req.GetExecs() {
		if exec.Id != "" {
			return nil, status.Error(codes.InvalidArgument, "Id must be empty")
		}
	}

	addedExecs, err := mongodb.AddExecToDb(ctx, req.GetExecs())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Execs{Execs: addedExecs}, nil
}

func (s *Server) GetExecs(ctx context.Context, req *pb.GetExecsRequest) (*pb.Execs, error) {
	// filtering
	filter, errs := buildFilter(req.Exec, &models.Exec{})
	if errs != nil {
		return nil, status.Error(codes.InvalidArgument, errs.Error())
	}

	// sorting
	sortOptions := buildSortOptions(req.GetSortBy())

	// access data
	execs, err := mongodb.GetExecsFromDB(ctx, sortOptions, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Execs{Execs: execs}, nil
}

func (s *Server) UpdateExecs(ctx context.Context, req *pb.Execs) (*pb.Execs, error) {
	updatedExecs, err := mongodb.ModifyExecsInDB(ctx, req.Execs)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Execs{Execs: updatedExecs}, nil
}

func (s *Server) DeleteExecs(ctx context.Context, req *pb.ExecIds) (*pb.DeleteExecsConfirmation, error) {
	ids := req.GetIds()

	var execIdsToDelete []string

	for _, exec := range ids {
		execIdsToDelete = append(execIdsToDelete, exec)
	}

	deletedIds, err := mongodb.DeleteExecsFromDB(ctx, execIdsToDelete)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteExecsConfirmation{Status: "success", DeletedIds: deletedIds}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.ExecLoginRequest) (*pb.ExecLoginResponse, error) {

	client, err := mongodb.CreateMongoClient()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer client.Disconnect(ctx)

	filter := bson.M{"username": req.GetUsername()}
	var exec models.Exec

	err = client.Database("school").Collection("execs").FindOne(ctx, filter).Decode(&exec)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrorHandler(err, "Exec not found")
		}
		return nil, utils.ErrorHandler(err, "Internal error")
	}

	if exec.InactiveStatus {
		return nil, status.Error(codes.PermissionDenied, "Exec is inactive")
	}

	err = utils.VerifyPassword(req.GetPassword(), exec.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid password")
	}

	token, err := utils.SignToken(exec.Id, exec.Username, exec.Role)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Could not generate token")
	}

	return &pb.ExecLoginResponse{
		Status: true,
		Token:  token,
	}, nil
}
