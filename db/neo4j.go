package db

import (
	"context"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var Driver neo4j.DriverWithContext

func InitDB(uri, username, password string) {
	var err error
	Driver, err = neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Fatalf("Failed to create driver: %v", err)
	}

	log.Println("Connected to Neo4j database")
	// createIndexes()
}

func CloseDB() {
	if Driver != nil {
		if err := Driver.Close(context.Background()); err != nil {
			log.Fatalf("Failed to close Neo4j driver: %v", err)
		}
	}
}
