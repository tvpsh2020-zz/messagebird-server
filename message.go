package main

import (
	"log"
	"sync"

	messagebird "github.com/messagebird/go-rest-api"
	"github.com/messagebird/go-rest-api/sms"
)

type RawMessage struct {
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

type messageQueue struct {
	sync.RWMutex
	List []Message
}

var MessageQueue = new(messageQueue)

func StoreMessageToQueue(rawMessage *RawMessage) error {

	// Validate or Error

	// var validator *Validator

	// if err := validator.validate(rawMessage); err != nil {

	// }

	messageBuilder := &MessageBuilder{
		RawMessage: rawMessage,
		Params: sms.Params{
			Type:       "binary",
			DataCoding: "",
		},
	}
	messages, err := messageBuilder.start()

	if err != nil {
		return err
	}

	MessageQueue.Lock()

	for _, message := range messages {
		MessageQueue.List = append(MessageQueue.List, message)
	}

	MessageQueue.Unlock()
	return nil
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
				log.Printf("Error: %#v\n", mbError)
			}
		}

		return
	}

	log.Printf("newMessage -> %#v ", newMessage)
}
