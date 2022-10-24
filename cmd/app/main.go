package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/firmfoundation/amqp-service/cmd/pkg/event"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	//select {}
	// start listening for messages
	log.Println("Listening for and consuming rabbitmq messages ..")

	//create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	//watch the queue and consume events
	err = consumer.Listen([]string{"auth"})
	if err != nil {
		panic(err)
	}

}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://rabbit:password@192.168.1.2:5672")
		if err != nil {
			fmt.Println("rabbitmq not yet ready ....")
			counts++
		} else {
			connection = c
			log.Println("Connected to rabbitmq service.")
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off ...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
