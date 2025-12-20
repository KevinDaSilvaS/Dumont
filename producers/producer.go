package producers

import (
	"context"
	"dumont/parser"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
}

func Connect() Producer {
	conn, _ := amqp.Dial("amqp://admin:admin@localhost:5672/")
	ch, _ := conn.Channel()
	q, _ := ch.QueueDeclare(
		"dumont", // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)

	return Producer{
		Connection: conn,
		Channel:    ch,
		Queue:      q,
	}
}

func (p Producer) Publish(transaction parser.Transaction) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, _ := json.Marshal(transaction)
	p.Channel.PublishWithContext(ctx, "", p.Queue.Name, false, false, amqp.Publishing{ContentType: "text/json", Body: []byte(body)})
}
