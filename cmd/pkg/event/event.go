package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExhange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", //name
		"topic",      //type
		true,         //durable?
		false,        // auto-deleted?
		false,        //internal?
		false,        // no-wait?
		nil,          // arguments?
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"queque-master", //name
		false,           // durable?
		false,           // delete when unsued?
		false,           // exclusive?
		false,           // no-wait?
		nil,             //arguments?
	)
}
