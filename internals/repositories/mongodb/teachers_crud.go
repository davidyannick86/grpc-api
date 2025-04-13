package mongodb

import (
	"context"

	"github.com/davidyannick86/grpc-api-mongodb/internals/models"
	"github.com/davidyannick86/grpc-api-mongodb/pkg/utils"
	pb "github.com/davidyannick86/grpc-api-mongodb/proto/gen"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetTeachersFromDB(ctx context.Context, sortOptions bson.D, filter bson.M) ([]*pb.Teacher, error) {
	client, err := CreateMongoClient()
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

	teachers, err := decodeEntities(ctx, cursor, func() *pb.Teacher { return &pb.Teacher{} }, func() *models.Teacher {
		return &models.Teacher{}
	})
	if err != nil {
		return nil, err
	}
	return teachers, nil
}

func AddTeacherToDb(ctx context.Context, teachersFomRequest []*pb.Teacher) ([]*pb.Teacher, error) {
	mongoClient, err := CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Failed to create MongoDB client")
	}
	defer mongoClient.Disconnect(ctx)

	newTeachers := make([]*models.Teacher, len(teachersFomRequest))

	for i, pbTeachereacher := range teachersFomRequest {
		newTeachers[i] = mapPbTeacherToModelTeacher(pbTeachereacher)
	}

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

		pbTeacher := mapModelTeacherToPb(*teacher)
		addedTeachers = append(addedTeachers, pbTeacher)
	}
	return addedTeachers, nil
}

func ModifyTeacherInDB(ctx context.Context, pbTeachers []*pb.Teacher) ([]*pb.Teacher, error) {
	client, err := CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Failed to create MongoDB client")
	}
	defer client.Disconnect(ctx)

	var updatedTeachers []*pb.Teacher

	for _, teacher := range pbTeachers {

		if teacher.Id == "" {
			return nil, utils.ErrorHandler(err, "Id must be set")
		}

		modelTeacher := mapPbTeacherToModelTeacher(teacher)
		objectID, err := primitive.ObjectIDFromHex(teacher.Id)

		if err != nil {
			return nil, utils.ErrorHandler(err, "Failed to convert ID to ObjectID")
		}

		// convert modelTeacher to bson document
		modelDoc, err := bson.Marshal(modelTeacher)
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

		_, err = client.Database("school").Collection("teachers").UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updatedDoc})
		if err != nil {
			return nil, utils.ErrorHandler(err, "Failed to update teacher in the database")
		}

		updatedTeacher := mapModelTeacherToPb(*modelTeacher)
		updatedTeachers = append(updatedTeachers, updatedTeacher)
	}
	return updatedTeachers, nil
}

func DeleteTeachersFromDB(ctx context.Context, teacherIdsToDelete []string) ([]string, error) {
	client, err := CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}
	defer client.Disconnect(ctx)

	objectIds := make([]primitive.ObjectID, len(teacherIdsToDelete))

	for i, id := range teacherIdsToDelete {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, utils.ErrorHandler(err, "invalid id")
		}
		objectIds[i] = objectId
	}

	filter := bson.M{"_id": bson.M{"$in": objectIds}}

	result, err := client.Database("school").Collection("teachers").DeleteMany(ctx, filter)
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}

	if result.DeletedCount == 0 {
		return nil, utils.ErrorHandler(err, "no teachers found")
	}

	deletedIds := make([]string, result.DeletedCount)
	for i, id := range objectIds {
		deletedIds[i] = id.Hex()
	}
	return deletedIds, nil
}
