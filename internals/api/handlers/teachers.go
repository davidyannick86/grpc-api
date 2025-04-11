package handlers

import (
	"context"
	"log"
	"reflect"
	"strings"

	"github.com/davidyannick86/grpc-api-mongodb/internals/models"
	"github.com/davidyannick86/grpc-api-mongodb/internals/repositories/mongodb"
	"github.com/davidyannick86/grpc-api-mongodb/pkg/utils"
	pb "github.com/davidyannick86/grpc-api-mongodb/proto/gen"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	filter, errs := buildFilterForTeacher(req)
	if errs != nil {
		return nil, utils.ErrorHandler(errs, "Invalid filter")
	}

	// sorting
	sortOptions := buildSortOptions(req.GetSortBy())

	// access data
	client, err := mongodb.CreateMongoClient()
	defer client.Disconnect(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	coll := client.Database("school").Collection("teachers")

	var cursor *mongo.Cursor

	if len(sortOptions) < 1 {
		cursor, err = coll.Find(ctx, filter)
	} else {
		cursor, err = coll.Find(ctx, filter, options.Find().SetSort(sortOptions))
	}

	defer cursor.Close(ctx)

	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal error")
	}

	//cursor, err = coll.Find(ctx, filter, options.Find().SetSort(sortOptions))

	var teachers []*pb.Teacher

	for cursor.Next(ctx) {
		var teacher models.Teacher
		err := cursor.Decode(&teacher)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Internal error")
		}
		teachers = append(teachers, &pb.Teacher{
			Id:        teacher.Id,
			FirstName: teacher.FirstName,
			LastName:  teacher.LastName,
			Email:     teacher.Email,
			Class:     teacher.Class,
			Subject:   teacher.Subject,
		})
	}

	return &pb.Teachers{Teachers: teachers}, nil
}

func buildFilterForTeacher(req *pb.GetTeachersRequest) (bson.M, error) {
	filter := bson.M{} // It's a map,  M is an unordered representation of a BSON document

	var modelTeacher models.Teacher

	modelVal := reflect.ValueOf(&modelTeacher).Elem()
	modelType := modelVal.Type()

	reqVal := reflect.ValueOf(req.Teacher).Elem()
	reqType := reqVal.Type()

	for i := range reqVal.NumField() {
		fieldVal := reqVal.Field(i)
		fieldName := reqType.Field(i).Name

		if fieldVal.IsValid() && !fieldVal.IsZero() {
			modelField := modelVal.FieldByName(fieldName)
			if modelField.IsValid() && modelField.CanSet() {
				modelField.Set(fieldVal)
			}
		}
	}

	// iterate over the modelTeacher to build filter using bson.M
	for i := range modelVal.NumField() {
		fieldVal := modelVal.Field(i)
		//fieldName := modelType.Field(i).Name

		if fieldVal.IsValid() && !fieldVal.IsZero() {
			bsonTag := modelType.Field(i).Tag.Get("bson")
			bsonTag = strings.TrimSuffix(bsonTag, ",omitempty")
			if bsonTag == "_id" {
				objectId, err := primitive.ObjectIDFromHex(req.Teacher.Id)
				if err != nil {
					return nil, utils.ErrorHandler(err, "Invalid ObjectId")
				}
				filter[bsonTag] = objectId
			} else {
				filter[bsonTag] = fieldVal.Interface().(string)
			}
		}
	}

	log.Println("filter", filter)

	return filter, nil
}

func buildSortOptions(sortFields []*pb.SortField) bson.D {
	var sortOptions bson.D

	for _, sortField := range sortFields {
		order := 1
		if sortField.GetOrder() == pb.Order_DESC {
			order = -1
		}
		sortOptions = append(sortOptions, bson.E{Key: sortField.Field, Value: order})
	}

	return sortOptions
}
