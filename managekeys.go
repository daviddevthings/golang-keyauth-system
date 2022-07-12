package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Key struct {
	Key            string `json:"key"`
	Timestamp      int64  `json:"timestamp"`
	User           int64  `json:"user"`
	MachineId      string `json:"machineid"`
	BoundToMachine bool   `json:"boundtomachine"`
}

func getAllKeys(ctx context.Context, coll *mongo.Collection) {
	cursor, currErr := coll.Find(ctx, bson.D{})

	if currErr != nil {
		panic(currErr)
	}
	defer cursor.Close(ctx)
	var keys []Key
	if err := cursor.All(ctx, &keys); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(keys)
}

func generateKey(coll *mongo.Collection) {
	rand.Seed(time.Now().UnixNano())
	var letters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	b := make([]byte, 24)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	newKey := Key{
		Key:            string(b),
		Timestamp:      time.Now().Unix(),
		User:           -1,
		MachineId:      "",
		BoundToMachine: false,
	}
	_, err := coll.InsertOne(context.TODO(), newKey)
	fmt.Printf("Added new key %s to the database", string(b))
	if err != nil {
		fmt.Println(err)
		return
	}
}
func deleteKey(ctx context.Context, coll *mongo.Collection, k string) {
	filter := bson.D{{Key: "key", Value: k}}
	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return
	}
	if result.DeletedCount == 0 {
		fmt.Printf("Key %s does not exist", k)
	} else {
		fmt.Printf("Deleted key %s from the database", k)
	}

}
func resetKey(ctx context.Context, coll *mongo.Collection, k string) {
	filter := bson.D{{Key: "key", Value: k}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "boundtomachine", Value: false}, {Key: "machineid", Value: ""}}}}

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return
	}
	if result.MatchedCount == 0 {
		fmt.Printf("Key %s does not exist", k)
	} else {
		fmt.Printf("Key %s was successfully reset", k)
	}

}
func bindKey(ctx context.Context, coll *mongo.Collection, k string, mID string) {
	filter := bson.D{{Key: "key", Value: k}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "boundtomachine", Value: true}, {Key: "machineid", Value: mID}}}}

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return
	}
	if result.MatchedCount == 0 {
		fmt.Printf("Key %s does not exist", k)
	} else {
		fmt.Printf("Key %s was successfully bound to machine %s", k, mID)
	}

}
func setUser(ctx context.Context, coll *mongo.Collection, k string, u int64) {
	filter := bson.D{{Key: "key", Value: k}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "user", Value: u}}}}

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return
	}
	if result.MatchedCount == 0 {
		fmt.Printf("Key %s does not exist", k)
	} else {
		fmt.Printf("User for key %s was successfully set", k)
	}

}
func authenticateKey(ctx context.Context, coll *mongo.Collection, k string, mID string) bool {
	filter := bson.D{{Key: "key", Value: k}}
	var result bson.D
	err := coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println("Key not found")
		return false
	}
	if result.Map()["machineid"] != mID {
		fmt.Println("Machine id doesn't matche one in the database")
		return false
	}
	return true
}
