package server

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func FetchDocuments(client *firestore.Client, ctx context.Context) ([]string, error) {
	iter := client.Collection("newsletter-emails").Documents(ctx)
	var recipients []string
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("failed to iterate: %v", err)
			return nil, err
		}

		email, ok := doc.Data()["email"].(string)
		if ok {
			recipients = append(recipients, email)
		} else {
			log.Fatalf("Invalid email format in document: %v", doc.Data())
		}

	}
	return recipients, nil
}
