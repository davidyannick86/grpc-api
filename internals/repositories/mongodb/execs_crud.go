package mongodb

import (
	"context"
	"time"

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

func AddExecToDb(ctx context.Context, execsFomRequest []*pb.Exec) ([]*pb.Exec, error) {
	mongoClient, err := CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Failed to create MongoDB client")
	}
	defer mongoClient.Disconnect(ctx)

	newExecs := make([]*models.Exec, len(execsFomRequest))

	for i, pbExec := range execsFomRequest {
		newExecs[i] = mapPbExecToModelExec(pbExec)
		hashedPassword, err := utils.HashPassword(newExecs[i].Password)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Failed to hash password")
		}
		newExecs[i].Password = hashedPassword
		currentTime := time.Now().Format(time.RFC3339)
		newExecs[i].UserCreatedAt = currentTime
		newExecs[i].InactiveStatus = false
	}

	var addedExecs []*pb.Exec

	for _, exec := range newExecs {
		result, err := mongoClient.Database("school").Collection("execs").InsertOne(ctx, exec)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Failed to insert exec into MongoDB")
		}

		objectId, ok := result.InsertedID.(primitive.ObjectID)
		if ok {
			exec.Id = objectId.Hex()
		}

		pbExec := mapModelExecToPb(*exec)
		addedExecs = append(addedExecs, pbExec)
	}
	return addedExecs, nil
}

func GetExecsFromDB(ctx context.Context, sortOptions bson.D, filter bson.M) ([]*pb.Exec, error) {
	client, err := CreateMongoClient()
	defer client.Disconnect(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	coll := client.Database("school").Collection("execs")

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

	execs, err := decodeEntities(ctx, cursor, func() *pb.Exec { return &pb.Exec{} }, func() *models.Exec {
		return &models.Exec{}
	})
	if err != nil {
		return nil, err
	}
	return execs, nil
}

func ModifyExecsInDB(ctx context.Context, pbExecs []*pb.Exec) ([]*pb.Exec, error) {
	client, err := CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Failed to create MongoDB client")
	}
	defer client.Disconnect(ctx)

	var updatedExecs []*pb.Exec

	for _, exec := range pbExecs {

		if exec.Id == "" {
			return nil, utils.ErrorHandler(err, "Id must be set")
		}

		modelExec := mapPbExecToModelExec(exec)
		objectID, err := primitive.ObjectIDFromHex(exec.Id)

		if err != nil {
			return nil, utils.ErrorHandler(err, "Failed to convert ID to ObjectID")
		}

		// convert modelExec to bson document
		modelDoc, err := bson.Marshal(modelExec)
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

		_, err = client.Database("school").Collection("execs").UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updatedDoc})
		if err != nil {
			return nil, utils.ErrorHandler(err, "Failed to update exec in the database")
		}

		updatedExec := mapModelExecToPb(*modelExec)
		updatedExecs = append(updatedExecs, updatedExec)
	}
	return updatedExecs, nil
}
