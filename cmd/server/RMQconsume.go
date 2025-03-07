package server

import "log"

func (s *server) Consume() {
	msgs, err := s.rabbitmq.Channel.Consume(
		"message_queue", // queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %v", msg.Body)
		}
	}()
}
