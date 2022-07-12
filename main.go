package main

import "fmt"

func main() {
	client, ctx := connectToDb("mongodb://localhost:27017")
	defer client.Disconnect(ctx)
	collection := client.Database("golangauth").Collection("keys")
	// generateKey(collection)
	// deleteKey(ctx, collection, "F94FNP75GNIMQ6GF978B84O9")
	// resetKey(ctx, collection, "F94FNP75GNIMQ6GF978B84O9")
	// setUser(ctx, collection, "S894V1LG5BM7DV82JV7Y7SIH", 448884623888351242)
	// getAllKeys(ctx, collection)
	fmt.Println(authenticateKey(ctx, collection, "S894V1LG5BM7DV82JV7Y7SIH", ""))
}
