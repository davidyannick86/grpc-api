package mongodb

import (
	"context"

	"github.com/davidyannick86/grpc-api-mongodb/internals/models"
	"github.com/davidyannick86/grpc-api-mongodb/pkg/utils"
	pb "github.com/davidyannick86/grpc-api-mongodb/proto/gen"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func GetStudentsFromDB(ctx context.Context, sortOptions bson.D, filter bson.M, pageNumber, pageSize uint32) ([]*pb.Student, error) {

	client, err := CreateMongoClient()
	defer client.Disconnect(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	coll := client.Database("school").Collection("students")

	findOptions := options.Find()
	findOptions.SetSkip((int64(pageNumber) - 1) * int64(pageSize))
	findOptions.SetLimit(int64(pageSize))

	if len(sortOptions) > 0 {
		findOptions.SetSort(sortOptions)
	}
	cursor, err := coll.Find(ctx, filter, findOptions)
	defer cursor.Close(ctx)

	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal error")
	}

	students, err := decodeEntities(ctx, cursor, func() *pb.Student { return &pb.Student{} }, func() *models.Student {
		return &models.Student{}
	})
	if err != nil {
		return nil, err
	}
	return students, nil
}

func ModifyStudentInDB(ctx context.Context, pbStudents []*pb.Student) ([]*pb.Student, error) {
	client, err := CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Failed to create MongoDB client")
	}
	defer client.Disconnect(ctx)

	var updatedStudents []*pb.Student

	for _, student := range pbStudents {

		if student.Id == "" {
			return nil, utils.ErrorHandler(err, "Id must be set")
		}

		modelStudent := mapPbStudentToModelStudent(student)
		objectID, err := primitive.ObjectIDFromHex(student.Id)

		if err != nil {
			return nil, utils.ErrorHandler(err, "Failed to convert ID to ObjectID")
		}

		// convert modelStudent to bson document
		modelDoc, err := bson.Marshal(modelStudent)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Failed to convert model to BSON")
		}

		// convert bson document to map
		var updatedDoc bson.M
		err = bson.Unmarshal(modelDoc, &updatedDoc)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Failed to convert BSON to map")
		}

		// remove the ID field from the updatedDoc
		delete(updatedDoc, "_id")

		_, err = client.Database("school").Collection("students").UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updatedDoc})
		if err != nil {
			return nil, utils.ErrorHandler(err, "Failed to update student in the database")
		}

		updatedStudent := mapModelStudentToPb(*modelStudent)
		updatedStudents = append(updatedStudents, updatedStudent)
	}
	return updatedStudents, nil
}

func DeleteStudentsFromDB(ctx context.Context, studentIdsToDelete []string) ([]string, error) {
	client, err := CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}
	defer client.Disconnect(ctx)

	objectIds := make([]primitive.ObjectID, len(studentIdsToDelete))

	for i, id := range studentIdsToDelete {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, utils.ErrorHandler(err, "invalid id")
		}
		objectIds[i] = objectId
	}

	filter := bson.M{"_id": bson.M{"$in": objectIds}}

	result, err := client.Database("school").Collection("students").DeleteMany(ctx, filter)
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}

	if result.DeletedCount == 0 {
		return nil, utils.ErrorHandler(err, "no students found")
	}

	deletedIds := make([]string, result.DeletedCount)
	for i, id := range objectIds {
		deletedIds[i] = id.Hex()
	}
	return deletedIds, nil
}
