package main

import (
	"encoding/json"
	"html"
	"log"
	"net/http"
	"strings"
)

type Route struct {
	Method  string
	URI     string
	Handler http.HandlerFunc
}

var routes = []Route{
	Route{"POST", "/api/message", sendMessageHandler},
	Route{"GET", "/api/balance", getBalanceHandler},
}

type apiResponse struct {
	Result string `json:"result"`
}

func initRouter() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		url := html.EscapeString(req.URL.Path)
		url = strings.TrimRight(url, "/")
		log.Println("[", req.Method, "]", url)

		res.Header().Set("Content-Type", "application/json")

		for _, route := range routes {
			if url == route.URI && req.Method == route.Method {

				route.Handler.ServeHTTP(res, req)

				return
			}

		}

		tmpResult := &apiResponse{
			Result: "Not found.",
		}

		res.WriteHeader(http.StatusNotFound)
		customRes, _ := json.Marshal(tmpResult)
		res.Write(customRes)
	})
}

func sendMessageHandler(res http.ResponseWriter, req *http.Request) {
	// 1. translate body
	var rawMessageBody RawMessageBody

	if err := json.NewDecoder(req.Body).Decode(&rawMessageBody); err != nil {
		http.Error(res, err.Error(), 400)
		log.Printf("%#v\n", err)
		return
	}

	// 2. validate all content

	// 3. store into message queue
	result, err := StoreMessageToQueue(&rawMessageBody)

	if err != nil {
		return
	}

	customRes, _ := json.Marshal(result)

	res.Write(customRes)
}

func getBalanceHandler(res http.ResponseWriter, req *http.Request) {
	balance := getBalance()
	balanceRes, _ := json.Marshal(balance)
	res.Write(balanceRes)
}
