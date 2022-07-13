package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func authEndpoint(ctx context.Context, collection *mongo.Collection) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		type authRequestBody struct {
			MachineId string `json:"machineid"`
			Key       string `json:"key"`
		}
		var requestBody authRequestBody

		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
		}
		if requestBody.Key == "" || requestBody.MachineId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "key or machine id missing",
			})
		} else {
			authResult := authenticateKey(ctx, collection, requestBody.Key, requestBody.MachineId)
			if authResult {
				c.JSON(http.StatusOK, gin.H{
					"message": "success",
				})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "not authenticated",
				})
			}
		}

	}

	return gin.HandlerFunc(fn)
}
func resetEndpoint(ctx context.Context, collection *mongo.Collection) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		key := c.Param("key")
		resetMessage, status := resetKey(ctx, collection, key)
		c.JSON(status, gin.H{
			"message": resetMessage,
		})

	}
	return fn
}
func deleteEndpoint(ctx context.Context, collection *mongo.Collection) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		key := c.Param("key")
		resetMessage, status := deleteKey(ctx, collection, key)
		c.JSON(status, gin.H{
			"message": resetMessage,
		})

	}
	return fn
}
func generateKeyEndpoint(collection *mongo.Collection) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		genkeymessage, status, errormessage := generateKey(collection)
		c.JSON(status, gin.H{
			"key":   genkeymessage,
			"error": errormessage,
		})

	}
	return fn
}
func setUserEndpoint(ctx context.Context, collection *mongo.Collection) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		type setUserRequestBody struct {
			User int64  `json:"user"`
			Key  string `json:"key"`
		}
		var requestBody setUserRequestBody

		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		} else if requestBody.Key == "" || requestBody.User == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "key or user id missing",
			})
		} else {
			message, status := setUser(ctx, collection, requestBody.Key, requestBody.User)

			c.IndentedJSON(status, gin.H{
				"message": message,
			})

		}

	}
	return fn
}
