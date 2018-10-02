package mock

import (
	"os"
	"testing"

	messagebird "github.com/messagebird/go-rest-api"
	"github.com/tvpsh2020/messagebird-server/config"
)

// StartMockServer can help you test your code
func StartMockServer(m *testing.M) {
	config.MBClient = messagebird.New("test_gshuPaZoeEG6ovbc8M79w0QyM")
	os.Exit(m.Run())
}
