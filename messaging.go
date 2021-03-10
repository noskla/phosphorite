package main

import (
	"encoding/json"
	"github.com/go-pg/pg/v10"
	"github.com/streadway/amqp"
	"log"
	"strconv"
)

var RabbitMQChannel *amqp.Channel
var RabbitMQRPCQueue amqp.Queue

func InitMessagingService(db *pg.DB) *amqp.Connection {
	var username = GetEnvVariable("PHO_MQ_USER", "guest")
	var password = GetEnvVariable("PHO_MQ_PASSWORD", "guest")
	var address = GetEnvVariable("PHO_MQ_ADDRESS", "127.0.0.1")
	var port = GetEnvVariable("PHO_MQ_PORT", "5672")
	var enabled = GetEnvVariable("PHO_MQ_ENABLED", "yes") == "yes"

	if !enabled {
		log.Println("RabbitMQ client disabled by environment variable.")
		return &amqp.Connection{}
	}

	conn, err := amqp.Dial("amqp://" + username + ":" + password + "@" + address + ":" + port + "/")
	if err != nil {
		log.Println("Error connecting with RabbitMQ service: ", err, "RPC disabled.")
		return conn
	}

	RabbitMQChannel, err = conn.Channel()
	if err != nil {
		log.Println("Error opening channel: ", err, "RPC disabled.")
		return conn
	}

	RabbitMQRPCQueue, err = RabbitMQChannel.QueueDeclare("phosphorite_rpc_queue",
		false, false, false, false, nil)

	err = RabbitMQChannel.Qos(1, 0, false)
	if err != nil {
		log.Println("Error setting QoS: ", err, "RPC disabled.")
		return conn
	}

	messages, err := RabbitMQChannel.Consume(RabbitMQRPCQueue.Name, "",
		false, false, false, false, nil)
	if err != nil {
		log.Println("Error registering consumer: ", err, "RPC disabled.")
		return conn
	}

	go MessagingServiceLoop(messages, db)
	log.Println("Messaging enabled")
	return conn

}

func MessagingServiceLoop(messages <-chan amqp.Delivery, db *pg.DB) {
	for msg := range messages {

		var rpccall map[string]interface{}
		if err := json.Unmarshal([]byte(msg.Body), &rpccall); err != nil {
			log.Println("Unhandled JSON RPC Call: ", err)
			continue
		}

		functionName, exist := rpccall["function_name"].(string)
		if !exist {
			continue
		}

		log.Println("RPC => " + functionName)
		var response = make(map[string]string)

		switch functionName {
		case "create_user":
			username := rpccall["username"].(string)
			password := rpccall["password"].(string)
			language := rpccall["language"].(string)
			userIP := rpccall["ip"].(string)

			err, code, userUUID := CreateUser(db, username, password, language, userIP)
			response["code"] = strconv.Itoa(code)
			if code != 1 {
				response["error"] = err.Error()
			} else {
				response["user_id"] = userUUID.String()
			}
		case "validate_user":
			username := rpccall["username"].(string)
			password := rpccall["password"].(string)
			saveDate := rpccall["save_date"].(string) == "yes"
			err, code, userUUID := ValidateUserPassword(db, username, password, saveDate)
			if code != 1 {
				response["error"] = err.Error()
			} else {
				response["user_id"] = userUUID.String()
			}

		}

		res, _ := json.Marshal(response)
		if err := RabbitMQChannel.Publish("", msg.ReplyTo, false, false, amqp.Publishing{
			ContentType: "application/json", CorrelationId: msg.CorrelationId, Body: []byte(res)}); err != nil {
			log.Println("Error publishing message: ", err)
			continue
		}
		_ = msg.Ack(false)

	}
}
