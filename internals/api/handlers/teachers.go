package handlers

import (
	"context"

	"github.com/davidyannick86/grpc-api-mongodb/internals/models"
	"github.com/davidyannick86/grpc-api-mongodb/internals/repositories/mongodb"
	pb "github.com/davidyannick86/grpc-api-mongodb/proto/gen"
	"go.mongodb.org/mongo-driver/mongo"
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

func (s *Server) GetTeachers(ctx context.Context, req *pb.GetTeachersRequest) (*pb.Teachers, error) {
	// filtering
	filter, errs := buildFilter(req.Teacher, &models.Teacher{})
	if errs != nil {
		return nil, status.Error(codes.InvalidArgument, errs.Error())
	}

	// sorting
	sortOptions := buildSortOptions(req.GetSortBy())

	// access data
	teachers, err := mongodb.GetTeachersFromDB(ctx, sortOptions, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Teachers{Teachers: teachers}, nil
}

func (s *Server) UpdateTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {
	updatedTeachers, err := mongodb.ModifyTeacherInDB(ctx, req.Teachers)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Teachers{Teachers: updatedTeachers}, nil
}

func (s *Server) DeleteTeachers(ctx context.Context, req *pb.TeacherIds) (*pb.DeleteTeachersConfirmation, error) {
	ids := req.GetIds()

	var teacherIdsToDelete []string

	for _, v := range ids {
		teacherIdsToDelete = append(teacherIdsToDelete, v.Id)
	}

	deletedIds, err := mongodb.DeleteTeachersFromDB(ctx, teacherIdsToDelete)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteTeachersConfirmation{Status: "success", DeletedIds: deletedIds}, nil
}

func (s *Server) GetStudentsByClassTeacher(ctx context.Context, req *pb.TeacherId) (*pb.Students, error) {

	teacherID := req.GetId()

	students, err := mongodb.GetStudentsByTeacherIDFromDB(ctx, teacherID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "No students found for this teacher")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Students{Students: students}, nil

}

func (s *Server) GetStudentCountByClassTeacher(ctx context.Context, req *pb.TeacherId) (*pb.StudentCount, error) {
	teacherID := req.GetId()

	count, err := mongodb.GetStudentCountByTeacherClass(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	return &pb.StudentCount{
		Status:       true,
		StudentCount: int32(count.StudentCount),
	}, nil

}
