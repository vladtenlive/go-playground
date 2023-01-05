package main

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func main() {
	ctx := context.Background()
	projectID := "test-project"

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
