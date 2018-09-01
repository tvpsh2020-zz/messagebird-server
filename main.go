package main

import (
	"net/http"
)

func init() {
	initConfig()
	initTaskManager()
	initRouter()
}

func main() {
	http.ListenAndServe(":"+serverPort, nil)
}
