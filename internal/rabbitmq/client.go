package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// Connecting to RabbitMQ and creating channel
func New(url string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{
		conn:    conn,
		channel: channel,
	}, nil
}

// Carefully close channel and connection, logging errors
func (c *Client) Close() error {
	if err := c.channel.Close(); err != nil {
		log.Printf("Failed to close RabbitMQ channel: %v", err)
	}
	if err := c.conn.Close(); err != nil {
		log.Printf("Failed to close RabbitMQ connection: %v", err)
	}
	return nil
}

// Createes queue (guaranteed) and publish message in it
func (c *Client) Publish(queue, message string) error {
	_, err := c.channel.QueueDeclare(
		queue,
		true,  // durable — очередь сохраняется при перезапуске сервера RabbitMQ
		false, // delete when unused — очередь не удаляется автоматически, когда к ней никто не подключён
		false, // exclusive — очередь не эксклюзивна для одного соединения, её могут использовать другие клиенты
		false, // no-wait — не ждать ответа от сервера после объявления очереди (false — ждать)
		nil,   // arguments — дополнительные аргументы, обычно nil
	)
	if err != nil {
		return err
	}

	return c.channel.Publish(
		"",    // exchange — пустая строка означает стандартный (default) exchange
		queue, // routing key — имя очереди, в которую отправляется сообщение
		false, // mandatory — если true и нет очереди для маршрутизации, сообщение не доставляется и возвращается отправителю; false — просто отбрасывается
		false, // immediate — если true и нет потребителя, сообщение не доставляется; false — ждать пока появится потребитель
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

func (c *Client) Consume(queue string) (<-chan amqp.Delivery, error) {
	_, err := c.channel.QueueDeclare(
		queue,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	return c.channel.Consume(
		queue,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}
