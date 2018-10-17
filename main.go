package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/ecojuntak/hastag-api/data"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"github.com/urfave/cli"
)

const QUEUE_NAME = "feed_system"

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func connectToServer() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	return conn
}

func openChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	return ch
}

func declareQueue(ch *amqp.Channel) amqp.Queue {
	q, err := ch.QueueDeclare(
		QUEUE_NAME, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")
	return q
}

func registerConsumer(ch *amqp.Channel, q amqp.Queue) (msgs <-chan amqp.Delivery) {
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	return
}

func startRabbitMQ() {
	conn := connectToServer()
	defer conn.Close()

	ch := openChannel(conn)
	defer ch.Close()

	q := declareQueue(ch)

	msgs := registerConsumer(ch, q)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			msg := string(d.Body[:])
			data.ProcessMessage(msg)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func AllHashtagHandler(w http.ResponseWriter, r *http.Request) {
	hashtags := data.GetAll()
	payload, _ := json.Marshal(hashtags)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
}

func SingleHashtagHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	feed_ids := data.GetFeedIds(name)
	json_str, _ := json.Marshal(feed_ids)

	queryParam := url.Values{"ids": {string(json_str[:])}}

	http.Get("http://localhost:8000/feeds/hashtag?" + queryParam.Encode())

	response, _ := http.Get("http://localhost:8000/feeds/hashtag?" + queryParam.Encode())

	payload, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err.Error)
	}

	fmt.Println(response.Body.Read)

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func startREST() {
	r := mux.NewRouter()
	r.HandleFunc("/hashtags", AllHashtagHandler)
	r.HandleFunc("/hashtags/{name}", SingleHashtagHandler)

	fmt.Println("REST server run on localhost:8088")
	log.Fatal(http.ListenAndServe(":8088", r))
}

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = "hashtag-api"
	cliApp.Version = "1.0.0"
	cliApp.Commands = []cli.Command{
		{
			Name:        "migrate",
			Description: "Run database migration",
			Action: func(c *cli.Context) error {
				err := data.RunMigration()
				if err != nil {
					log.Println(err)
				}
				return err
			},
		},
		{
			Name:        "start-amqp",
			Description: "Start listening to RabbitMQ Server",
			Action: func(c *cli.Context) error {
				startRabbitMQ()
				return nil
			},
		},
		{
			Name:        "start-rest",
			Description: "Start listening to REST API",
			Action: func(c *cli.Context) error {
				startREST()
				return nil
			},
		},
	}

	cliApp.Run(os.Args)
}
