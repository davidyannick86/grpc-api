package handlers

import (
	"context"
	"fmt"
	"reflect"

	"github.com/davidyannick86/grpc-api-mongodb/internals/models"
	"github.com/davidyannick86/grpc-api-mongodb/internals/repositories/mongodb"
	"github.com/davidyannick86/grpc-api-mongodb/pkg/utils"
	pb "github.com/davidyannick86/grpc-api-mongodb/proto/gen"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Server) AddTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {
	mongoClient, err := mongodb.CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Failed to create MongoDB client")
	}
	defer mongoClient.Disconnect(ctx)

	newTeachers := make([]*models.Teacher, len(req.GetTeachers()))

	for i, pbTeachereacher := range req.GetTeachers() {
		modelTeacher := models.Teacher{FirstName: "Tagada"}
		pbVal := reflect.ValueOf(pbTeachereacher).Elem()
		modelVal := reflect.ValueOf(&modelTeacher).Elem()

		for i := range pbVal.NumField() {
			pbField := pbVal.Field(i)
			fieldName := pbVal.Type().Field(i).Name

			modelField := modelVal.FieldByName(fieldName)
			if modelField.IsValid() && modelField.CanSet() {
				modelField.Set(pbField)
			}
		}
		newTeachers[i] = &modelTeacher
	}

	fmt.Println("newTeachers", newTeachers)

	// // return &pb.Teachers{Teachers: addedTeacher}, nil
	var addedTeachers []*pb.Teacher

	for _, teacher := range newTeachers {
		result, err := mongoClient.Database("school").Collection("teachers").InsertOne(ctx, teacher)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Failed to insert teacher into MongoDB")
		}

		objectId, ok := result.InsertedID.(primitive.ObjectID)
		if ok {
			teacher.Id = objectId.Hex()
		}

		pbTeacher := &pb.Teacher{}
		modelVal := reflect.ValueOf(*teacher)
		pbVal := reflect.ValueOf(pbTeacher).Elem() // Utiliser Elem() pour obtenir la valeur pointée

		for i := range modelVal.NumField() { // Correction de la boucle
			modelField := modelVal.Field(i)

			modelFieldType := modelVal.Type().Field(i)

			pbField := pbVal.FieldByName(modelFieldType.Name)

			if pbField.IsValid() && pbField.CanSet() {
				pbField.Set(modelField)
			}
		}
		addedTeachers = append(addedTeachers, pbTeacher)
	}

	return &pb.Teachers{Teachers: addedTeachers}, nil
}
