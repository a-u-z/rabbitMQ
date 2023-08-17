package main

import (
	"log"

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

	// 建立通道，在一個與 rabbitMQ 的連線上，可以有多個通道，所以可以併發建立 conn.Channel，達到更多的吞吐量
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 宣告隊列
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable 簡單的範例，不需要持久化，durable:true 的原因是當 rabbitMQ 崩潰了，重新啟動後，隊列的訊息還在，如果訊息沒有那麼重要，那麼可以設定成 false
		false,   // delete when unused，队列在不再被使用（没有消费者连接到队列或是没有绑定的交换机，滿足其一）时会被自动删除。
		false,   // exclusive，只允许创建它的连接（即通常是同一个连接上的不同信道）访问队列，用於临时队列與限制訪問
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 宣告消費者
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer，空字串表示由 rabbitMQ 產生隨機的消費者，需要這個欄位原因是在需要更確定知道訊息被哪個消費者消費了
		true,   // auto-ack，自動發送 ack 給 rabbitMQ，false 的話需要額外的程式碼來執行如「第四講」
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	// 消費
	go func() {
		for d := range msgs { // 不斷的從 msgs 這個通道拿取東西，當通道關閉或是沒有訊息時，會自動退出
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// var forever chan struct{}
// 	<-forever
// 用於堵塞，讓 func 不結束
