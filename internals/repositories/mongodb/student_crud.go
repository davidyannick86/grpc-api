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

	students, err := decodeEntities(ctx, cursor, func() *pb.Student { return &pb.Student{} }, newModel)
	if err != nil {
		return nil, err
	}
	return students, nil
}
