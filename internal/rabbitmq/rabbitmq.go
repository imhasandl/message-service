package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

const (
	// ExchangeName is the name of the topic exchange used for notifications.
	ExchangeName = "notifications.topic"
	// QueueName is the name of the queue for the notification service.
	QueueName    = "notification_service_queue"
	// RoutingKey is the routing key used to bind the queue to the exchange for notification messages.
	RoutingKey   = "message-service.notification"
)

// RabbitMQ encapsulates the RabbitMQ connection and channel.
type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// NewRabbitMQ creates a new RabbitMQ instance and establishes a connection and channel.
// It takes the RabbitMQ server URL as input and returns a pointer to the RabbitMQ struct or an error if connection fails.
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("can't connect to rabbit mq: %v", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("can't connect to the channel: %v", err)
		return nil, err
	}

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}, nil
}

// Close cleanly closes the RabbitMQ channel and connection.
// It logs any errors encountered during the closing process.
func (r *RabbitMQ) Close() {
	if r.Channel != nil {
		if err := r.Channel.Close(); err != nil {
			log.Printf("error closing channel: %v", err)
		}
	}
	if r.Conn != nil {
		if err := r.Conn.Close(); err != nil {
			log.Printf("error closing connection: %v", err)
		}
	}
}
