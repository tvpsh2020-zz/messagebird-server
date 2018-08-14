package main

import (
	"fmt"

	"github.com/tvpsh2020/messagebird-server/lib"

	messagebird "github.com/messagebird/go-rest-api"
)

func GetBalance() *messagebird.Balance {
	balance, err := lib.MBClient.Balance()

	if err != nil {
		// messagebird.ErrResponse means custom JSON errors.
		if err == messagebird.ErrResponse {
			for _, mbError := range balance.Errors {
				fmt.Printf("Error: %#v\n", mbError)
			}

			return nil
		}

		fmt.Println(err)
		return nil
	}

	fmt.Println("MBClient payment -> ", balance.Payment)
	fmt.Println("MBClient type -> ", balance.Type)
	fmt.Println("MBClient amount -> ", balance.Amount)

	return balance
}
