package main

import (
	"fmt"
	"log"
	"time"

	"github.com/davidyannick86/grpc-api-mongodb/internals/repositories/mongodb"
)

func main() {

	log.Println("Seeding database...")
	start := time.Now()
	defer func() {
		log.Printf("Database seeded in %s", time.Since(start))
	}()

	err := mongodb.SeedDatabase()
	if err != nil {
		fmt.Println("Error seeding database:", err)
		return
	}
	log.Println("Database seeded successfully.")

}
