package rabbit_test

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/robotjoosen/go-rabbit"
	"github.com/wagslane/go-rabbitmq"
)

func ExampleConsumer() {
	conn, err := rabbit.NewConnection("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic(err)
	}

	c, err := rabbit.NewConsumer(
		conn,
		"test-exchange",
		"test-routing-key",
		"test-queue-name",
	)
	if err != nil {
		panic(err)
	}

	if err = c.Run(func(d rabbitmq.Delivery) (action rabbitmq.Action) {
		return rabbitmq.Ack
	}); err != nil {
		panic(err)
	}

	c.Close()
}

func ExampleRunConsumer() {
	ctx, cnl := context.WithCancel(context.Background())

	conn, err := rabbit.NewConnection("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic(err)
	}

	go func() {
		if err = rabbit.RunConsumer(
			conn,
			"test-exchange",
			"test-routing-key",
			"test-queue-name",
			func(d rabbitmq.Delivery) (action rabbitmq.Action) {
				return rabbitmq.Ack
			},
			ctx,
		); err != nil {
			panic(err)
		}
	}()

	time.Sleep(10 * time.Second)

	cnl() // stop it from running
}

func ExamplePublisher() {
	conn, err := rabbit.NewConnection("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic(err)
	}

	p, err := rabbit.NewPublisher(conn, "test-exchange")
	if err != nil {
		panic(err)
	}

	if err = rabbit.Publish(
		[]byte("hello"),
		"test-routing-key",
		"test-exchange", uuid.New().String(),
		p,
	); err != nil {
		panic(err)
	}
}
