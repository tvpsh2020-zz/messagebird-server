package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tvpsh2020/messagebird-server/lib"
)

/*
TODO

1. receive msg from my api
2. body check: >160 (or unicode > 70) should be split (how to let server know this is a split SMS?)
3. body check: cannot be empty
4. body check: invalid word
5. maximum split SMS is 255
6. MessageBird API can only accept 1 req per second, let all received messages be handle in queue
7. refactor into OOP
8. make a http route to handle flexible uri pattern
9. make config loader easy to add-on

*/

func init() {
	lib.AppInitialize()
}

func main() {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/api/message", sendMessageHandler)
	http.HandleFunc("/api/balance", getBalanceHandler)
	http.HandleFunc("/api/test", apiTestHandler)
	http.ListenAndServe(":"+lib.GlobalConfig.Port, nil)
}

type APIResponse struct {
	Result string `json:"result"`
}

func defaultHandler(res http.ResponseWriter, req *http.Request) {
	tmpResult := &APIResponse{
		Result: "Hi there.",
	}

	res.Header().Set("Content-Type", "application/json")
	customRes, _ := json.Marshal(tmpResult)
	res.Write(customRes)
}

type SendMessageBody struct {
	Recipients string
	Originator string
	Body       string
}

func sendMessageHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("H: ", req.Header)

	if req.Header["Content-Type"][0] == "application/json" {
		fmt.Println("Good~")
	}

	if req.Method == "POST" {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}

		var msg SendMessageBody
		err = json.Unmarshal(b, &msg)

		fmt.Println(msg)

		result := SendSMSToMessageBird(&msg)

		customRes, _ := json.Marshal(result)
		res.Header().Set("Content-Type", "application/json")

		res.Write(customRes)
	} else {
		tmpResult := &APIResponse{
			Result: "Bad request.",
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusBadRequest)
		customRes, _ := json.Marshal(tmpResult)
		res.Write(customRes)
	}
}

func getBalanceHandler(res http.ResponseWriter, req *http.Request) {
	balance := GetBalance()

	balanceRes, _ := json.Marshal(balance)

	res.Header().Set("Content-Type", "application/json")
	res.Write(balanceRes)
}

func apiTestHandler(res http.ResponseWriter, req *http.Request) {
	message := SMSSendTest()

	messageRes, _ := json.Marshal(message)

	res.Header().Set("Content-Type", "application/json")
	res.Write(messageRes)
}
