package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// 做連線
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// 宣告通道
	ch, err := conn.Channel() // 建立通道，在一個與 rabbitMQ 的連線上，可以有多個通道，所以可以併發建立 conn.Channel，達到更多的吞吐量
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 宣告隊列
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 發送消息
	body := "Hello Worldd!"
	err = ch.PublishWithContext(ctx,
		"",     // exchange，空字串表示默認交換機，如果在消費者端沒有使用交換機，那這邊就是空字串
		q.Name, // routing key，看要發送到哪個隊列
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}
