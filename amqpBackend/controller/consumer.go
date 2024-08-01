package main

import (
	"fmt"

	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("Application Connsumer")
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

	msgs, err := ch.Consume(
		"testAPI",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fmt.Printf("Receive Message: %s\n", d.Body)
		}

	}()
	fmt.Println("Success receive message ")
	fmt.Println("[*] - Wait for message")
	<-forever
}
