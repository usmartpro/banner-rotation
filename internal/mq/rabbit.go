package mq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/usmartpro/banner-rotation/internal/app"
)

type Rabbit struct {
	exchange    string
	queue       string
	consumerTag string
	channel     *amqp.Channel
	logger      app.Logger
}

func (q *Rabbit) GetChannel() *amqp.Channel {
	return q.channel
}

func NewRabbit(
	ctx context.Context,
	dsn string,
	exchange string,
	queue string,
	logger app.Logger,
) (*Rabbit, error) {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("error connect to rabbit (%s): %w", dsn, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error open rabbit channel (%s): %w", dsn, err)
	}

	if len(exchange) > 0 {
		if err = ch.ExchangeDeclare(
			exchange,
			amqp.ExchangeDirect,
			true,
			false,
			false,
			false,
			nil,
		); err != nil {
			return nil, fmt.Errorf("error declare exchange (%s): %w", exchange, err)
		}
	}

	q, err := ch.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error declare queue (%s): %w", queue, err)
	}

	if err = ch.QueueBind(
		q.Name,
		q.Name,
		exchange,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("error bind queue: %w", err)
	}

	go func() {
		<-ctx.Done()
		ch.Close()
		conn.Close()
	}()

	return &Rabbit{
		exchange:    exchange,
		queue:       queue,
		consumerTag: "banners-consumer",
		channel:     ch,
		logger:      logger,
	}, nil
}

func (q *Rabbit) Add(n app.StatEvent) error {
	body, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("error marshall StatEvent: %w", err)
	}

	if err = q.channel.Publish(
		q.exchange,
		q.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		}); err != nil {
		return fmt.Errorf("error publish StatEvent: %w", err)
	}

	return nil
}

func (q *Rabbit) GetNotificationChannel() (<-chan app.StatEvent, error) {
	ch := make(chan app.StatEvent)

	deliveries, err := q.channel.Consume(
		q.queue,
		q.consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error consume queue (%s): %w", q.queue, err)
	}

	go func() {
		for d := range deliveries {
			var notification app.StatEvent
			if err := json.Unmarshal(d.Body, &notification); err != nil {
				q.logger.Error("error unmarshal notification: %s", err)
				continue
			}

			ch <- notification

			d.Ack(false)
		}

		close(ch)
	}()

	return ch, nil
}
