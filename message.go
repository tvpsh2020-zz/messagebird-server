package main

import (
	"log"
	"sync"

	messagebird "github.com/messagebird/go-rest-api"
	"github.com/messagebird/go-rest-api/sms"
)

// IRawMessage is an interface for raw message
type IRawMessage struct {
	Recipients string
	Originator string
	Body       string
}

// IMessage is an interface for a single message
type IMessage struct {
	Originator string
	Body       string
	Recipients []string
	Params     sms.Params
}

// IMessageQueue is an interface for message queue
type IMessageQueue struct {
	sync.RWMutex
	List []IMessage
}

// MessageQueue is global access message queue instance
var MessageQueue = new(IMessageQueue)

func StoreMessageToQueue(rawMessage *IRawMessage) error {
	messageBuilder := &IMessageBuilder{
		IRawMessage: rawMessage,
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

func SendSMSToMessageBirdV2(message IMessage) {

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
