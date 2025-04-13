package mongodb

import (
	"context"

	"github.com/davidyannick86/grpc-api-mongodb/internals/models"
	"github.com/davidyannick86/grpc-api-mongodb/pkg/utils"
	pb "github.com/davidyannick86/grpc-api-mongodb/proto/gen"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
