package mongodb

import (
	"context"

	"github.com/davidyannick86/grpc-api-mongodb/internals/models"
	"github.com/davidyannick86/grpc-api-mongodb/pkg/utils"
	pb "github.com/davidyannick86/grpc-api-mongodb/proto/gen"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddStudentToDb(ctx context.Context, studentsFomRequest []*pb.Student) ([]*pb.Student, error) {
	mongoClient, err := CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Failed to create MongoDB client")
	}
	defer mongoClient.Disconnect(ctx)

	newStudents := make([]*models.Student, len(studentsFomRequest))

	for i, pbStudent := range studentsFomRequest {
		newStudents[i] = mapPbStudentToModelStudent(pbStudent)
	}

	var addedStudents []*pb.Student

	for _, student := range newStudents {
		result, err := mongoClient.Database("school").Collection("students").InsertOne(ctx, student)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Failed to insert student into MongoDB")
		}

		objectId, ok := result.InsertedID.(primitive.ObjectID)
		if ok {
			student.Id = objectId.Hex()
		}

		pbStudent := mapModelStudentToPb(*student)
		addedStudents = append(addedStudents, pbStudent)
	}
	return addedStudents, nil
}
