package main

import (
	"fmt"

	messagebird "github.com/messagebird/go-rest-api"
)

func getBalance() *messagebird.Balance {
	balance, err := mbClient.Balance()

	if err != nil {
		if err == messagebird.ErrResponse {
			for _, mbError := range balance.Errors {
				fmt.Printf("Error: %#v\n", mbError)
			}

			return nil
		}

		fmt.Println(err)
		return nil
	}

	return balance
}
