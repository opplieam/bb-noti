package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("Error connecting to nats server: ", err)
	}
	defer nc.Close()

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        "update",
		Description: "Message for update",
		Subjects: []string{
			"update.>",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	i := 0
	for {
		i++
		time.Sleep(1 * time.Second)
		categoryName := fmt.Sprintf("category-%d", i)
		_, err = js.Publish(ctx, fmt.Sprintf("update.category.%d", i), []byte(categoryName))
		if err != nil {
			log.Println("Error publishing message: ", err)
			continue
		}
		log.Printf("Published Order: [%d]", i)
	}
}
