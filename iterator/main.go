package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func main() {
	ctx := context.Background()
	projectID := "test-project"

	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8084")
	}

	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	defer firestoreClient.Close()

	iter := firestoreClient.Collection("Items").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err != nil && errors.Is(err, iterator.Done) {
			return
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(doc.Data())
	}
}
