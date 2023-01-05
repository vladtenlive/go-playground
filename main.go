package main

import (
	"context"
	"log"
	"time"

	pb "go-playground/proto"

	"cloud.google.com/go/firestore"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	projectID := "test-project"

	logger := watermill.NewStdLogger(false, false)
	subscriber, err := googlecloud.NewSubscriber(
		googlecloud.SubscriberConfig{
			GenerateSubscriptionName: func(topic string) string {
				return "test-sub_" + topic
			},
			ProjectID: projectID,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	// firerstore

	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	defer firestoreClient.Close()

	//

	messages, err := subscriber.Subscribe(context.Background(), "example.topic")
	if err != nil {
		panic(err)
	}

	go process(messages, firestoreClient)

	publisher, err := googlecloud.NewPublisher(googlecloud.PublisherConfig{
		ProjectID: projectID,
	}, logger)
	if err != nil {
		panic(err)
	}

	publishMessages(publisher)
}

func publishMessages(publisher message.Publisher) {
	for {
		item := pb.Item{
			Id:      uuid.New().String(),
			Payload: "DATA: " + uuid.New().String(),
		}
		data, err := proto.Marshal(&item)
		if err != nil {
			panic(err)
		}

		msg := message.NewMessage(watermill.NewUUID(), data)

		if err := publisher.Publish("example.topic", msg); err != nil {
			panic(err)
		}

		time.Sleep(time.Second)
	}
}

func process(messages <-chan *message.Message, firestoreClient *firestore.Client) {
	for msg := range messages {
		log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))

		item := pb.Item{}
		err := proto.Unmarshal(msg.Payload, &item)
		if err != nil {
			panic(err)
		}

		_, err = firestoreClient.Doc("Items/"+msg.UUID).Set(context.Background(), map[string]string{
			"id":      item.Id,
			"payload": item.Payload,
		})
		if err != nil {
			panic(err)
		}

		msg.Ack()
	}
}
