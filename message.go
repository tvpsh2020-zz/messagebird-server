package main

import (
	"fmt"

	"github.com/tvpsh2020/messagebird-server/lib"

	messagebird "github.com/messagebird/go-rest-api"
)

func SMSSendTest() *messagebird.Message {
	params := &messagebird.MessageParams{Reference: "MyReference"}

	message, err := lib.MBClient.NewMessage(
		"Jimmy",
		[]string{"886931077193"},
		"Hi There.",
		params)

	if err != nil {
		if err == messagebird.ErrResponse {
			for _, mbError := range message.Errors {
				fmt.Printf("Error: %#v\n", mbError)
			}

			return nil
		}

		fmt.Println(err)
		return nil
	}

	fmt.Println(message)

	return message
}

func SendSMSToMessageBird(sms *SendMessageBody) *messagebird.Message {
	params := &messagebird.MessageParams{Reference: "MyReference"}

	fmt.Println(sms.Originator)
	fmt.Println(sms.Recipients)
	fmt.Println(sms.Body)

	message, err := lib.MBClient.NewMessage(
		sms.Originator,
		[]string{sms.Recipients},
		sms.Body,
		params)

	if err != nil {
		if err == messagebird.ErrResponse {
			for _, mbError := range message.Errors {
				fmt.Printf("Error: %#v\n", mbError)
			}

			return nil
		}

		fmt.Println(err)
		return nil
	}

	fmt.Println(message)

	return message
}
