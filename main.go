package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	client, ctx := connectToDb("mongodb://localhost:27017")
	defer client.Disconnect(ctx)
	collection := client.Database("golangauth").Collection("keys")
	r := gin.Default()
	r.GET("/allkeys", func(c *gin.Context) {
		c.JSON(http.StatusOK, getAllKeys(ctx, collection))
	})
	r.POST("/auth", authEndpoint(ctx, collection))
	r.GET("/reset/:key", resetEndpoint(ctx, collection))
	r.GET("/delete/:key", deleteEndpoint(ctx, collection))
	r.GET("/genkey", generateKeyEndpoint(collection))
	r.POST("/setuser", setUserEndpoint(ctx, collection))
	r.Run()
}
