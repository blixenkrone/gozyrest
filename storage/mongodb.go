package storage

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DatabaseRef dynamic referencing db and collection in mongo?
type DatabaseRef struct {
	Database   string
	Collection string
}

// ObjectID -
type ObjectID struct {
	ObjectID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

// Tip -
type Tip struct {
	ObjectID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
}

var ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)

// InsertOneItem abstraction of the real mongo DB method insertOne()
func (dbRef *DatabaseRef) InsertOneItem(document interface{}) (*mongo.InsertOneResult, error) {
	col, err := openMongo(dbRef)
	if err != nil {
		fmt.Printf("Error opening coll: %s", err)
		return nil, err
	}
	res, err := col.InsertOne(ctx, document)
	if err != nil {
		fmt.Printf("Error inserting: %s\n", err)
		return nil, err
	}
	return res, nil
}

// GetAllTips gets a single item from MongoDB
func (dbRef *DatabaseRef) GetAllTips(ch chan<- []*Tip, wg *sync.WaitGroup) {
	wg.Add(1)
	var tips []*Tip
	col, err := openMongo(dbRef)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error: %s", err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var tip Tip
		cursor.Decode(&tip)
		tips = append(tips, &tip)
	}
	if err := cursor.Err(); err != nil {
		log.Printf("Error: %s", err)
	}
	ch <- tips
	close(ch)
	wg.Done()
}

// openMongo create a connection from client to database
func openMongo(dbRef *DatabaseRef) (*mongo.Collection, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	// err = client.Ping(ctx, readpref.Primary())
	// if err != nil {
	// 	return nil, err
	// }
	collection := client.Database(dbRef.Database).Collection(dbRef.Collection)
	return collection, nil
}
