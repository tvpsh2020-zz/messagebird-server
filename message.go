package main

import (
	"fmt"
	"sync"

	messagebird "github.com/messagebird/go-rest-api"
)

type RawMessageBody struct {
	Recipients string
	Originator string
	Body       string
}

type Message struct {
	Originator string
	Body       string
	Recipients []string
	Params     messagebird.MessageParams
}

// var MessageQueue []Message

type messageQueue struct {
	sync.RWMutex
	List []Message
}

var MessageQueue = new(messageQueue)

func StoreMessageToQueue(message *RawMessageBody) (*apiResponse, error) {
	// split to an array

	messageBuilder := &MessageBuilder{RawMessageBody: message}
	messages := messageBuilder.start()

	// send to queue,

	MessageQueue.Lock()

	for _, message := range messages {
		MessageQueue.List = append(MessageQueue.List, message)
	}

	MessageQueue.Unlock()

	tmpResult := &apiResponse{
		Result: "OK.",
	}

	return tmpResult, nil
}

func SendSMSToMessageBirdV2(message Message) {

	newMessage, err := mbClient.NewMessage(
		message.Originator,
		message.Recipients,
		message.Body,
		&message.Params)

	if err != nil {
		if err == messagebird.ErrResponse {
			for _, mbError := range newMessage.Errors {
				fmt.Printf("Error: %#v\n", mbError)
			}
		}

		fmt.Println(err)
	}

}
