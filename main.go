package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	client, ctx := connectToDb("mongodb://localhost:27017")
	defer client.Disconnect(ctx)
	collection := client.Database("golangauth").Collection("keys")
	r := gin.Default()
	r.GET("/allkeys", AllKeysEndpoint(ctx, collection))
	r.POST("/auth", AuthEndpoint(ctx, collection))
	r.GET("/reset/:key", ResetEndpoint(ctx, collection))
	r.GET("/delete/:key", DeleteEndpoint(ctx, collection))
	r.GET("/genkey", GenerateKeyEndpoint(collection))
	r.POST("/setuser", SetUserEndpoint(ctx, collection))
	r.Run()
}
