package mongodb

import (
	"context"
	"reflect"

	"github.com/davidyannick86/grpc-api-mongodb/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func decodeEntities[T any, M any](ctx context.Context, cursor *mongo.Cursor, newEntity func() *T, newModel func() *M) ([]*T, error) {
	var entities []*T

	for cursor.Next(ctx) {
		model := newModel()
		err := cursor.Decode(&model)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Internal error")
		}

		entity := newEntity()
		modelVal := reflect.ValueOf(model).Elem() // Utiliser Elem() pour obtenir la valeur pointées
		pbVal := reflect.ValueOf(entity).Elem()   // Utiliser Elem() pour obtenir la valeur pointées

		for i := range modelVal.NumField() {
			modelField := modelVal.Field(i)
			modelFieldName := modelVal.Type().Field(i).Name

			pbField := pbVal.FieldByName(modelFieldName)
			if pbField.IsValid() && pbField.CanSet() {
				pbField.Set(modelField)
			}
		}

		entities = append(entities, entity)
	}

	err := cursor.Err()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal error")
	}

	return entities, nil
}
