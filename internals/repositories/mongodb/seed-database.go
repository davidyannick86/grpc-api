package mongodb

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/davidyannick86/grpc-api-mongodb/internals/models"
	"github.com/davidyannick86/grpc-api-mongodb/pkg/utils"
)

func SeedDatabase() error {
	client, err := CreateMongoClient()
	if err != nil {
		log.Printf("Error: %v", utils.ErrorHandler(err, "Failed to create MongoDB client"))
		return err
	}
	defer client.Disconnect(context.Background())

	// Seed Teachers
	teachersJsonPath := "./json/teachersdata.json"
	teachersData, err := os.ReadFile(teachersJsonPath)
	if err != nil {
		log.Printf("Error: %v", utils.ErrorHandler(err, "Failed to read teachers JSON file"))
		return err
	}

	teachers := []*models.Teacher{}
	err = json.Unmarshal(teachersData, &teachers)
	if err != nil {
		log.Printf("Error: %v", utils.ErrorHandler(err, "Failed to unmarshal teachers JSON"))
		return err
	}
	err = client.Database("school").Collection("teachers").Drop(context.Background())
	if err != nil {
		log.Printf("Error: %v", utils.ErrorHandler(err, "Failed to drop teachers collection"))
		return err
	}

	for _, teacher := range teachers {
		_, err = client.Database("school").Collection("teachers").InsertOne(context.Background(), teacher)
		if err != nil {
			log.Printf("Error: %v", utils.ErrorHandler(err, "Failed to insert teacher"))
			continue
		}
	}

	// Seed Students
	studentsJsonPath := "./json/studentsdata.json"
	studentsData, err := os.ReadFile(studentsJsonPath)
	if err != nil {
		log.Printf("Error: %v", utils.ErrorHandler(err, "Failed to read students JSON file"))
		return err
	}

	students := []*models.Student{}
	err = json.Unmarshal(studentsData, &students)
	if err != nil {
		log.Printf("Error: %v", utils.ErrorHandler(err, "Failed to unmarshal students JSON"))
		return err
	}

	err = client.Database("school").Collection("students").Drop(context.Background())
	if err != nil {
		log.Printf("Error: %v", utils.ErrorHandler(err, "Failed to drop students collection"))
		return err
	}

	for _, student := range students {
		_, err = client.Database("school").Collection("students").InsertOne(context.Background(), student)
		if err != nil {
			log.Printf("Error: %v", utils.ErrorHandler(err, "Failed to insert student"))
			continue
		}
	}
	// Seed Execs
	execsJsonPath := "./json/execsdata.json"
	execsData, err := os.ReadFile(execsJsonPath)
	if err != nil {
		log.Printf("Error: %v", utils.ErrorHandler(err, "Failed to read execs JSON file"))
		return err
	}

	execs := []*models.Exec{}
	err = json.Unmarshal(execsData, &execs)
	if err != nil {
		log.Printf("Error: %v", utils.ErrorHandler(err, "Failed to unmarshal execs JSON"))
		return err
	}
	err = client.Database("school").Collection("execs").Drop(context.Background())
	if err != nil {
		log.Printf("Error: %v", utils.ErrorHandler(err, "Failed to drop execs collection"))
		return err
	}
	for _, exec := range execs {
		hashedPassword, err := utils.HashPassword(exec.Password)
		if err != nil {
			log.Printf("Error: %v", utils.ErrorHandler(err, "Failed to hash password"))
			continue
		}
		exec.Password = hashedPassword
		exec.UserCreatedAt = time.Now().Format(time.RFC3339)
		exec.InactiveStatus = false
		_, err = client.Database("school").Collection("execs").InsertOne(context.Background(), exec)
		if err != nil {
			log.Printf("Error: %v", utils.ErrorHandler(err, "Failed to insert exec"))
			continue
		}
	}
	return nil
}
