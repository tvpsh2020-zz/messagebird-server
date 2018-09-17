package main

import (
	"encoding/json"
	"html"
	"log"
	"net/http"
	"strings"
)

type route struct {
	Method  string
	URI     string
	Handler http.HandlerFunc
}

var routes = []route{
	route{"POST", "/api/message", sendMessageHandler},
	route{"GET", "/api/balance", getBalanceHandler},
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
	var rawMessage *IRawMessage

	if err := json.NewDecoder(req.Body).Decode(&rawMessage); err != nil {
		http.Error(res, err.Error(), 400)
		log.Printf("%#v\n", err)
		return
	}

	if err := StoreMessageToQueue(rawMessage); err != nil {
		log.Println()
		tmpResult := &apiResponse{
			Result: err.Error(),
		}

		customRes, _ := json.Marshal(tmpResult)
		res.Write(customRes)
	} else {
		tmpResult := &apiResponse{
			Result: "OK",
		}

		customRes, _ := json.Marshal(tmpResult)
		res.Write(customRes)
	}

}

func getBalanceHandler(res http.ResponseWriter, req *http.Request) {
	balance := getBalance()
	balanceRes, _ := json.Marshal(balance)
	res.Write(balanceRes)
}
