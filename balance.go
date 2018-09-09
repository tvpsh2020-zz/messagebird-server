package main

import (
	"fmt"

	messagebird "github.com/messagebird/go-rest-api"
	"github.com/messagebird/go-rest-api/balance"
)

func getBalance() *balance.Balance {
	balance, err := balance.Read(mbClient)

	if err != nil {
		switch errResp := err.(type) {
		case messagebird.ErrorResponse:
			for _, mbError := range errResp.Errors {
				fmt.Printf("Error: %#v\n", mbError)
			}
		}

		return nil
	}

	return balance
}
