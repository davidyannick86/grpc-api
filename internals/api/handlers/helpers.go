package handlers

import (
	"reflect"
	"strings"

	"github.com/davidyannick86/grpc-api-mongodb/pkg/utils"
	pb "github.com/davidyannick86/grpc-api-mongodb/proto/gen"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func buildFilter(object any, model any) (bson.M, error) {
	filter := bson.M{} // It's a map,  M is an unordered representation of a BSON document

	if object == nil || reflect.ValueOf(object).IsNil() {
		return filter, nil
	}

	modelVal := reflect.ValueOf(model).Elem()
	modelType := modelVal.Type()

	reqVal := reflect.ValueOf(object).Elem()
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
		fieldName := modelType.Field(i).Name

		if fieldVal.IsValid() && !fieldVal.IsZero() {
			bsonTag := modelType.Field(i).Tag.Get("bson")
			bsonTag = strings.TrimSuffix(bsonTag, ",omitempty")
			if bsonTag == "_id" {
				objectId, err := primitive.ObjectIDFromHex(reqVal.FieldByName(fieldName).Interface().(string))
				if err != nil {
					return nil, utils.ErrorHandler(err, "Invalid ObjectId")
				}
				filter[bsonTag] = objectId
			} else {
				filter[bsonTag] = fieldVal.Interface().(string)
			}
		}
	}

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
