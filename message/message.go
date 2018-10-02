package message

import (
	"log"
	"sync"

	"github.com/messagebird/go-rest-api/sms"
	"github.com/tvpsh2020/messagebird-server/config"
	"github.com/tvpsh2020/messagebird-server/consts"
)

// IRawMessage is an interface for raw message
type IRawMessage struct {
	Recipients string `json:"recipients"`
	Originator string `json:"originator"`
	Body       string `json:"body"`
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

// Queue for store messages
var Queue = new(IMessageQueue)

func sendToMessageBird(message IMessage) (*sms.Message, error) {
	newMessage, err := sms.Create(
		config.MBClient,
		message.Originator,
		message.Recipients,
		message.Body,
		&message.Params)

	if err != nil {
		log.Printf("error from MessageBird: %#v", err)

		return nil, err
	}

	log.Printf("newMessage: %#v ", newMessage)
	return newMessage, nil
}

// StoreToMessageQueue will keep your message into queue
func StoreToMessageQueue(rawMessage *IRawMessage) error {
	messageBuilder := &IMessageBuilder{
		IRawMessage:       rawMessage,
		ShouldConcatenate: false,
		Params: sms.Params{
			Type:       consts.SMSType,
			DataCoding: "",
		},
	}

	messages, err := messageBuilder.start()

	if err != nil {
		return err
	}

	Queue.Lock()

	for _, message := range messages {
		Queue.List = append(Queue.List, message)
	}

	Queue.Unlock()

	return nil
}

// SendFromMessageQueue will bring your message to the world
func SendFromMessageQueue() {
	Queue.RLock()

	if len(Queue.List) > 0 {
		sendToMessageBird(Queue.List[0])
		Queue.RUnlock()
		Queue.Lock()
		Queue.List = Queue.List[1:]
		Queue.Unlock()
	} else {
		Queue.RUnlock()
	}
}
