package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Key struct {
	Key       string `json:"key"`
	Timestamp int64  `json:"timestamp"`
	User      int64  `json:"user"`
	MachineId string `json:"machineid"`
}

func GetAllKeys(ctx context.Context, coll *mongo.Collection) ([]Key, int) {
	cursor, currErr := coll.Find(ctx, bson.D{})

	if currErr != nil {
		panic(currErr)
	}
	defer cursor.Close(ctx)
	var keys []Key
	if err := cursor.All(ctx, &keys); err != nil {
		return []Key{}, http.StatusInternalServerError
	}
	return keys, http.StatusOK
}

func GenerateKey(coll *mongo.Collection) (Key, int, string) {
	rand.Seed(time.Now().UnixNano())
	var letters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	b := make([]byte, 24)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	newKey := Key{
		Key:       string(b),
		Timestamp: time.Now().UnixNano(),
		User:      -1,
		MachineId: "none",
	}
	_, err := coll.InsertOne(context.TODO(), newKey)
	if err != nil {
		return Key{}, http.StatusInternalServerError, "error generating key"
	}
	return newKey, http.StatusOK, "none"
}
func DeleteKey(ctx context.Context, coll *mongo.Collection, k string) (string, int) {
	filter := bson.D{{Key: "key", Value: k}}
	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {

		return "Error deleting key", http.StatusInternalServerError
	}
	if result.DeletedCount == 0 {
		return fmt.Sprintf("Key %s does not exist", k), http.StatusNotFound
	} else {
		return fmt.Sprintf("Deleted key %s from the database", k), http.StatusOK
	}

}
func ResetKey(ctx context.Context, coll *mongo.Collection, k string) (string, int) {
	filter := bson.D{{Key: "key", Value: k}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "machineid", Value: "none"}}}}

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {

		return "error resetting key", http.StatusInternalServerError
	}
	if result.MatchedCount == 0 {
		return fmt.Sprintf("Key %s does not exist", k), http.StatusNotFound
	} else {
		return fmt.Sprintf("Key %s was successfully reset", k), http.StatusOK
	}

}
func SetUser(ctx context.Context, coll *mongo.Collection, k string, u int64) (string, int) {
	filter := bson.D{{Key: "key", Value: k}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "user", Value: u}}}}

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return "Error setting user", http.StatusInternalServerError
	}
	if result.MatchedCount == 0 {
		return fmt.Sprintf("Key %s does not exist", k), http.StatusNotFound
	} else {
		return fmt.Sprintf("User for key %s was successfully set", k), http.StatusOK
	}

}
func AuthenticateKey(ctx context.Context, coll *mongo.Collection, k string, mID string) bool {
	filter := bson.D{{Key: "key", Value: k}}
	var result bson.D
	err := coll.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		//key not found
		return false
	}
	if result.Map()["machineid"] == "none" {
		filter := bson.D{{Key: "key", Value: k}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "machineid", Value: mID}}}}
		_, err := coll.UpdateOne(ctx, filter, update)
		if err != nil {
			//error updating
			return false
		}
		return true
	} else if result.Map()["machineid"] != mID {
		//Machine id doesn't matche one in the database
		return false
	}
	return true
}
