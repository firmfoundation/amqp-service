package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn:      conn,
		queueName: "queque-master",
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExhange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		err := ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range messages {
			var payload AuthPayload
			json.Unmarshal(d.Body, &payload)
			go handlePayload(payload)
		}
	}()

	fmt.Printf("waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever

	return nil
}

func handlePayload(payload AuthPayload) {
	switch payload.Email {
	case "log", "event":
		//log whatever we get
		// err := logEvent(payload)
		// if err != nil {
		// 	log.Print(err)
		// }

	case "auth":
		//authenticate
		// fmt.Println("YELLLLLLLLLLLLLLLLLLLL")
		// logAuth(payload)
	default:
		fmt.Println("YELLLLLLLLLLLLLLLLLLLL")
		go logAuth(payload)
		// err := logEvent(payload)
		// if err != nil {
		// 	log.Print(err)
		// }
	}
}

// func logEvent(entry Payload) error {
// 	return nil
// }

func logAuth(entry AuthPayload) error {
	fmt.Println("arrived ...")
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	// call the service
	requestURL := fmt.Sprintf("http://192.168.1.2:%d/authenticate", 8081)
	request, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(jsonData))
	//request, err := http.NewRequest("POST", "http://auth-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	//response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("***************", err)
		return err
	}
	defer response.Body.Close()

	fmt.Println("greeeeeeeeen", response.Body)
	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		return err
	} else if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
