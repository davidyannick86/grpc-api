package handlers

import (
	"context"
	"fmt"
	"reflect"

	"github.com/davidyannick86/grpc-api-mongodb/internals/models"
	"github.com/davidyannick86/grpc-api-mongodb/internals/repositories/mongodb"
	"github.com/davidyannick86/grpc-api-mongodb/pkg/utils"
	pb "github.com/davidyannick86/grpc-api-mongodb/proto/gen"
)

func (s *Server) AddTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {
	mongoClient, err := mongodb.CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Failed to create MongoDB client")
	}
	defer mongoClient.Disconnect(ctx)

	newTeachers := make([]*models.Teacher, len(req.GetTeachers()))

	for i, t := range req.GetTeachers() {
		newTeachers[i] = mapPbTeacherToModelTeacher(t)
	}

	for _, teacher := range newTeachers {
		fmt.Println(teacher.FirstName)
		fmt.Println(teacher.LastName)
	}

	return nil, nil

}

func mapPbTeacherToModelTeacher(pbTeacher *pb.Teacher) *models.Teacher {
	modelTeacher := models.Teacher{}
	pbVal := reflect.ValueOf(pbTeacher).Elem()
	modelVal := reflect.ValueOf(&modelTeacher).Elem()

	for i := range pbVal.NumField() {
		pbField := pbVal.Field(i)
		fieldName := pbVal.Type().Field(i).Name

		modelField := modelVal.FieldByName(fieldName)
		if modelField.IsValid() && modelField.CanSet() {
			modelField.Set(pbField)
		}
	}
	return &modelTeacher
}
