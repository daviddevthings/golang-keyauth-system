package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthEndpoint(ctx context.Context, collection *mongo.Collection) gin.HandlerFunc {
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
			authResult := AuthenticateKey(ctx, collection, requestBody.Key, requestBody.MachineId)
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
func ResetEndpoint(ctx context.Context, collection *mongo.Collection) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		key := c.Param("key")
		resetMessage, status := ResetKey(ctx, collection, key)
		c.JSON(status, gin.H{
			"message": resetMessage,
		})

	}
	return fn
}
func DeleteEndpoint(ctx context.Context, collection *mongo.Collection) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		key := c.Param("key")
		resetMessage, status := DeleteKey(ctx, collection, key)
		c.JSON(status, gin.H{
			"message": resetMessage,
		})

	}
	return fn
}
func GenerateKeyEndpoint(collection *mongo.Collection) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		genkeymessage, status, errormessage := GenerateKey(collection)
		c.JSON(status, gin.H{
			"key":   genkeymessage,
			"error": errormessage,
		})

	}
	return fn
}
func SetUserEndpoint(ctx context.Context, collection *mongo.Collection) gin.HandlerFunc {
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
			message, status := SetUser(ctx, collection, requestBody.Key, requestBody.User)

			c.IndentedJSON(status, gin.H{
				"message": message,
			})

		}

	}
	return fn
}
func AllKeysEndpoint(ctx context.Context, collection *mongo.Collection) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		res, status := GetAllKeys(ctx, collection)
		c.JSON(status, res)
	}
	return fn
}
