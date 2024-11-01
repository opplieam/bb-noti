package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/opplieam/bb-noti-api/internal/category"
	"github.com/opplieam/bb-noti-api/internal/state"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	clientState := state.NewClientState()
	// -----------------------
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("Error connecting to nats server: ", err)
	}
	defer nc.Close()

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal("Error connecting to JetStream server: ", err)
	}
	ctx := context.Background()
	stream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        "update",
		Description: "Message for update",
		Subjects: []string{
			"update.>",
		},
	})
	if err != nil {
		log.Fatal("Error creating stream: ", err)
	}
	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:    "update_processor",
		Durable: "update_processor",
	})
	if err != nil {
		log.Fatal("Error creating update_processor: ", err)
	}
	conCtx, err := consumer.Consume(func(msg jetstream.Msg) {
		_ = msg.Ack()
		for _, ch := range clientState.GetAllClients() {
			ch <- fmt.Sprintf("%s : %s", msg.Subject(), msg.Data())
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conCtx.Stop()
	// ----------------------

	var r *gin.Engine
	r = gin.Default()
	r.Use(gin.Recovery())

	categoryHandler := category.NewHandler(clientState)
	r.GET("/category", categoryHandler.SSE)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r.Handler(),
	}

	serverErrors := make(chan error, 1)
	go func() {
		defer close(serverErrors)
		serverErrors <- srv.ListenAndServe()
	}()

	select {
	case err = <-serverErrors:
		return fmt.Errorf("server error: %s", err)
	case <-shutdown:
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			_ = srv.Close()
			return fmt.Errorf("could not shutdown gratefuly: %w", err)
		}
	}
	return nil
}
