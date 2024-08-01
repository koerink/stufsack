package main

import (
	"amqpBackend/gnutils"
	"fmt"
)

func main() {
	fmt.Println("Go RabbitMQ")

	message := map[string]interface{}{
		"Type": "transaction",
		"Body": map[string]interface{}{
			"username":      "data1",
			"password":      "somepassword",
			"Email":         "koerink1@gmail.com",
			"Password":      "somepassword",
			"IsVerifyEmail": "false",
			"Profile": map[string]string{
				"Bio":     "-",
				"Address": "-",
				"Parrent": "-",
			},
			"CourseId": []int{1},
		},
	}
	gnutils.PublishMessage("exchDirectTrx", "exchange", "direct", message, "registerTrx")
}
