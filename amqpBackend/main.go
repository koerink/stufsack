package main

import (
	"fmt"

	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("Go RabbitMQ")
	conn, err := amqp.Dial("amqp://application:apps1234@apimasagi.my.id:5672/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()

	fmt.Print("Success Connected")

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer ch.Close()

	ttl := int32(10000)
	queue, err := ch.QueueDeclare("testAPI", true, false, false, false, amqp.Table{"x-message-ttl": ttl})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(queue)

	err = ch.Publish(
		"",
		"testAPI",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Hello Wordl"),
		},
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("Success Publish message to queu")
}
