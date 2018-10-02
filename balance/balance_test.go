package balance

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/tvpsh2020/messagebird-server/mock"
)

func TestMain(m *testing.M) {
	mock.StartMockServer(m)
	os.Exit(m.Run())
}
func TestShouldGetBalanceCorrectly(t *testing.T) {
	balance, err := GetBalance()

	if err != nil {
		t.Errorf("Expected error to be nil but got %#v", err)
	}

	expectedResponse := `{"Payment":"prepaid","Type":"credits","Amount":10.1}`
	response, _ := json.Marshal(balance)
	actualResponse := string(response)

	if expectedResponse != actualResponse {
		t.Fatalf("Expected response to be %s but got %s", expectedResponse, actualResponse)
	}
}
