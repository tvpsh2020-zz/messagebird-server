package main

import (
	"fmt"
	"sync"

	messagebird "github.com/messagebird/go-rest-api"
	"github.com/messagebird/go-rest-api/sms"
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
	Params     sms.Params
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

	newMessage, err := sms.Create(
		mbClient,
		message.Originator,
		message.Recipients,
		message.Body,
		&message.Params)

	if err != nil {
		switch errResp := err.(type) {
		case messagebird.ErrorResponse:
			for _, mbError := range errResp.Errors {
				fmt.Printf("Error: %#v\n", mbError)
			}
		}

		return
	}

	fmt.Println(newMessage)
}
