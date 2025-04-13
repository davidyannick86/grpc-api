package mongodb

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/davidyannick86/grpc-api-mongodb/internals/models"
	"github.com/davidyannick86/grpc-api-mongodb/pkg/utils"

	pb "github.com/davidyannick86/grpc-api-mongodb/proto/gen"
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

func mapModelToPb[M any, P any](model M, newPb func() *P) *P {

	pbStruct := newPb()
	modelVal := reflect.ValueOf(model)
	pbVal := reflect.ValueOf(pbStruct).Elem() // Utiliser Elem() pour obtenir la valeur pointée

	for i := 0; i < modelVal.NumField(); i++ { // Correction de la boucle
		modelField := modelVal.Field(i)

		modelFieldType := modelVal.Type().Field(i)

		pbField := pbVal.FieldByName(modelFieldType.Name)

		if pbField.IsValid() && pbField.CanSet() {
			pbField.Set(modelField)
		}
	}
	return pbStruct
}

func mapModelTeacherToPb(teacherModel models.Teacher) *pb.Teacher {
	return mapModelToPb(teacherModel, func() *pb.Teacher {
		return &pb.Teacher{}
	})
}

func mapModelStudentToPb(studentModel models.Student) *pb.Student {
	return mapModelToPb(studentModel, func() *pb.Student {
		return &pb.Student{}
	})
}

func mapModelExecToPb(execModel models.Exec) *pb.Exec {
	return mapModelToPb(execModel, func() *pb.Exec {
		return &pb.Exec{}
	})
}

func mapPbToModel[M any, P any](pbStruct P, newModel func() *M) *M {

	modelStuct := newModel()
	pbVal := reflect.ValueOf(pbStruct).Elem()
	modelVal := reflect.ValueOf(modelStuct).Elem()

	for i := 0; i < pbVal.NumField(); i++ {
		pbField := pbVal.Field(i)
		fieldName := pbVal.Type().Field(i).Name

		modelField := modelVal.FieldByName(fieldName)
		if modelField.IsValid() && modelField.CanSet() {
			modelField.Set(pbField)
		}
	}
	return modelStuct
}

func mapPbTeacherToModelTeacher(pbTeacher *pb.Teacher) *models.Teacher {
	return mapPbToModel(pbTeacher, func() *models.Teacher {
		return &models.Teacher{}
	})
}

func mapPbStudentToModelStudent(pbStudent *pb.Student) *models.Student {
	return mapPbToModel(pbStudent, func() *models.Student {
		return &models.Student{}
	})
}

func mapPbExecToModelExec(pbExec *pb.Exec) *models.Exec {
	return mapPbToModel(pbExec, func() *models.Exec {
		return &models.Exec{}
	})
}
