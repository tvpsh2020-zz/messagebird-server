package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/tvpsh2020/messagebird-server/message"
	"github.com/tvpsh2020/messagebird-server/mock"
)

func TestMain(m *testing.M) {
	mock.StartMockServer(m)
	os.Exit(m.Run())
}

func TestShouldSendMessageCorrectly(t *testing.T) {
	rawMessage := message.IRawMessage{
		Recipients: "886987654321",
		Originator: "Jimmy",
		Body:       "Hello",
	}

	body, _ := json.Marshal(rawMessage)

	request, err := http.NewRequest("POST", "/message", bytes.NewReader(body))

	if err != nil {
		t.Errorf("Expected error to be nil but got %#v", err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(sendMessageHandler)
	handler.ServeHTTP(recorder, request)

	expectedStatusCode := http.StatusCreated
	actualStatusCode := recorder.Code

	if expectedStatusCode != actualStatusCode {
		t.Errorf("Expected http status code to be %d but got %d", expectedStatusCode, actualStatusCode)
	}

	var actualResponse apiResponse
	err = json.NewDecoder(recorder.Body).Decode(&actualResponse)

	expectedResponse := apiResponse{
		Result: "success",
		Data:   map[string]string{"message": "your message will be send if no bad things happen, please check server log"},
	}

	if expectedResponse.Result != actualResponse.Result {
		t.Errorf("Expected api response result to be %v but got %v", expectedResponse.Result, actualResponse.Result)
	}

}
