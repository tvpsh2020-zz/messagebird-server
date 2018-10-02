package router

import (
	"encoding/json"
	"html"
	"log"
	"net/http"
	"strings"

	"github.com/tvpsh2020/messagebird-server/balance"
	"github.com/tvpsh2020/messagebird-server/consts"
	"github.com/tvpsh2020/messagebird-server/message"
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
	Result string      `json:"result"`
	Data   interface{} `json:"data"`
}

func (a apiResponse) send() []byte {
	apiResponse := apiResponse{
		Result: a.Result,
		Data:   a.Data,
	}

	result, _ := json.Marshal(apiResponse)

	return result
}

// Initialize router
func Initialize() {
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

		apiResponse := apiResponse{
			Result: consts.APIFail,
			Data:   map[string]string{"message": "url not found"},
		}

		res.WriteHeader(http.StatusNotFound)
		customRes, _ := json.Marshal(apiResponse)
		res.Write(customRes)
	})
}

func sendMessageHandler(res http.ResponseWriter, req *http.Request) {
	rawMessage := &message.IRawMessage{}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(rawMessage); err != nil {
		apiResponse := apiResponse{
			Result: consts.APIFail,
			Data:   map[string]string{"message": err.Error()},
		}

		log.Printf("error from server: %#v", err.Error())
		res.WriteHeader(http.StatusBadRequest)
		res.Write(apiResponse.send())
		return
	}

	if err := message.StoreToMessageQueue(rawMessage); err != nil {
		apiResponse := apiResponse{
			Result: consts.APIFail,
			Data:   map[string]string{"message": err.Error()},
		}

		log.Printf("error from server: %#v", err.Error())
		res.WriteHeader(http.StatusBadRequest)
		res.Write(apiResponse.send())
		return
	}

	apiResponse := apiResponse{
		Result: consts.APISuccess,
		Data:   map[string]string{"message": "your message will be send if no bad things happen, please check server log"},
	}

	res.WriteHeader(http.StatusCreated)
	res.Write(apiResponse.send())
}

func getBalanceHandler(res http.ResponseWriter, req *http.Request) {
	balance, err := balance.GetBalance()

	if err != nil {
		apiResponse := apiResponse{
			Result: consts.APIFail,
			Data:   map[string]string{"message": err.Error()},
		}

		log.Printf("error from MessageBird: %#v", err)
		res.WriteHeader(http.StatusBadRequest)
		res.Write(apiResponse.send())
		return
	}

	apiResponse := apiResponse{
		Result: consts.APISuccess,
		Data:   balance,
	}

	res.Write(apiResponse.send())
}
