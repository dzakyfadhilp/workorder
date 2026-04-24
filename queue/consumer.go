package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"workorder-api/model"
	"workorder-api/repository"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	workorderRepo *repository.WorkorderRepository
}

func NewConsumer(conn *amqp.Connection, workorderRepo *repository.WorkorderRepository) (*Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Set QoS - process 1 message at a time
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	return &Consumer{
		conn:          conn,
		channel:       ch,
		workorderRepo: workorderRepo,
	}, nil
}

func (c *Consumer) Start() error {
	msgs, err := c.channel.Consume(
		QueueName, // queue
		"",        // consumer tag
		false,     // auto-ack (manual ack for reliability)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("Consumer started, waiting for messages...")

	// Process messages
	go func() {
		for msg := range msgs {
			c.processMessage(msg)
		}
	}()

	return nil
}

func (c *Consumer) processMessage(delivery amqp.Delivery) {
	var msg Message
	if err := json.Unmarshal(delivery.Body, &msg); err != nil {
		log.Printf("[%s] ERROR: Failed to unmarshal message: %v", "UNKNOWN", err)
		delivery.Nack(false, false) // Don't requeue invalid messages
		return
	}

	log.Printf("[%s] Processing message: function=%s", msg.RequestID, msg.Function)

	// Route based on function
	var err error
	switch msg.Function {
	case "ff_updateWorkorder":
		err = c.processUpdateWorkorder(&msg)
	default:
		log.Printf("[%s] WARN: Unknown function: %s", msg.RequestID, msg.Function)
		delivery.Ack(false) // Ack to remove from queue
		return
	}

	if err != nil {
		log.Printf("[%s] ERROR: Failed to process message: %v", msg.RequestID, err)
		// Nack with requeue - will retry
		delivery.Nack(false, true)
		return
	}

	log.Printf("[%s] Message processed successfully", msg.RequestID)
	delivery.Ack(false)
}

func (c *Consumer) processUpdateWorkorder(msg *Message) error {
	var payload model.WorkorderRequest
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return fmt.Errorf("failed to parse payload: %w", err)
	}

	// Validate
	if payload.Req.Wonum == "" || payload.Req.Status == "" || payload.Req.Siteid == "" {
		return fmt.Errorf("validation failed: wonum, status, siteid are required")
	}

	// Save to database
	if err := c.workorderRepo.UpsertWorkorder(&payload, msg.RequestID); err != nil {
		return fmt.Errorf("failed to save workorder: %w", err)
	}

	log.Printf("[%s] Workorder saved: wonum=%s, status=%s", msg.RequestID, payload.Req.Wonum, payload.Req.Status)
	return nil
}

func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
}
