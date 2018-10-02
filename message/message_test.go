package message

import (
	"os"
	"testing"
	"unicode/utf8"

	"github.com/messagebird/go-rest-api/sms"

	"github.com/tvpsh2020/messagebird-server/consts"
	"github.com/tvpsh2020/messagebird-server/mock"
)

func TestMain(m *testing.M) {
	mock.StartMockServer(m)
	os.Exit(m.Run())
}

func TestShouldSendMessageToMessageBirdCorrectly(t *testing.T) {
	param := sms.Params{
		Reference: "Jimmy",
	}

	message := IMessage{
		Originator: "Jimmy",
		Body:       "Hello",
		Recipients: []string{"886987654321"},
		Params:     param,
	}

	newMessage, err := sendToMessageBird(message)

	if err != nil {
		t.Errorf("Expected error to be nil but got %#v", err)
	}

	expectedBody := message.Body
	expectedOriginator := message.Originator
	actualBody := newMessage.Body
	actualOriginator := newMessage.Originator

	if expectedBody != actualBody {
		t.Errorf("Expected body to be %s but got %s", expectedBody, actualBody)
	}

	if expectedOriginator != actualOriginator {
		t.Errorf("Expected body to be %s but got %s", expectedOriginator, actualOriginator)
	}
}

func TestShouldValidateOriginatorCorrectly(t *testing.T) {
	rawMessage := &IRawMessage{
		Recipients: "886987654321",
		Originator: "Jimmy",
		Body:       "Hello",
	}

	messageBuilder := &IMessageBuilder{
		IRawMessage:       rawMessage,
		ShouldConcatenate: false,
		Params: sms.Params{
			Type:       consts.SMSType,
			DataCoding: "",
		},
	}

	if err := messageBuilder.validateOriginator(); err != nil {
		t.Errorf("Expected error to be nil but got %#v", err)
	}
}

func TestShouldFixRecipientsCorrectly(t *testing.T) {
	rawMessage := &IRawMessage{
		Recipients: "886987654321, +886987654321,  , ++abcd987654321, j7k6j5h4bv3n",
		Originator: "Jimmy",
		Body:       "Hello",
	}

	messageBuilder := &IMessageBuilder{
		IRawMessage:       rawMessage,
		ShouldConcatenate: false,
		Params: sms.Params{
			Type:       consts.SMSType,
			DataCoding: "",
		},
	}

	expectedFixedRecipients := "886987654321,886987654321,,987654321,76543"
	actualFixesRecipients := messageBuilder.fixRecipients()

	if expectedFixedRecipients != actualFixesRecipients {
		t.Errorf("Expected fixed recipients to be %s but got %s", expectedFixedRecipients, actualFixesRecipients)
	}
}

func TestShouldSplitRecipientsCorrectly(t *testing.T) {
	rawMessage := &IRawMessage{
		Recipients: "886987654321",
		Originator: "Jimmy",
		Body:       "Hello",
	}

	messageBuilder := &IMessageBuilder{
		IRawMessage:       rawMessage,
		ShouldConcatenate: false,
		Params: sms.Params{
			Type:       consts.SMSType,
			DataCoding: "",
		},
	}

	if err := messageBuilder.splitRecipients(); err != nil {
		t.Errorf("Expected error to be nil but got %#v", err)
	}
}

func TestShouldCheckBodyDataCodingCorrectly(t *testing.T) {
	rawMessage := &IRawMessage{
		Recipients: "886987654321",
		Originator: "Jimmy",
		Body:       "Hello",
	}

	messageBuilder := &IMessageBuilder{
		IRawMessage:       rawMessage,
		ShouldConcatenate: false,
		Params: sms.Params{
			Type:       consts.SMSType,
			DataCoding: "",
		},
	}

	if err := messageBuilder.checkBodyDataCoding(); err != nil {
		t.Errorf("Expected error to be nil but got %#v", err)
	}
}

func TestShouldCountBodyLengthCorrectly(t *testing.T) {
	rawMessage := &IRawMessage{
		Recipients: "886987654321",
		Originator: "Jimmy",
		Body:       "Hello",
	}

	messageBuilder := &IMessageBuilder{
		IRawMessage:       rawMessage,
		ShouldConcatenate: false,
		Params: sms.Params{
			Type:       consts.SMSType,
			DataCoding: "plain",
		},
	}

	if err := messageBuilder.countBodyLength(); err != nil {
		t.Errorf("Expected error to be nil but got %#v", err)
	}
}

func TestShouldMakeSingleMessageCorrectly(t *testing.T) {
	rawMessage := &IRawMessage{
		Recipients: "886987654321",
		Originator: "Jimmy",
		Body:       "Hello",
	}

	messageBuilder := &IMessageBuilder{
		IRawMessage:       rawMessage,
		ShouldConcatenate: false,
		Params: sms.Params{
			Type:       consts.SMSType,
			DataCoding: "plain",
		},
	}

	expectedResult := IMessage{
		Originator: rawMessage.Originator,
		Body:       rawMessage.Body,
		Recipients: []string{rawMessage.Recipients},
		Params:     messageBuilder.Params,
	}

	actualResult := messageBuilder.makeSingleMessage(rawMessage.Body)

	if expectedResult.Originator != actualResult.Originator {
		t.Errorf("Expected originator to be %s but got %s", expectedResult.Originator, actualResult.Originator)
	}

	if expectedResult.Body != actualResult.Body {
		t.Errorf("Expected body to be %s but got %s", expectedResult.Body, actualResult.Body)
	}
}

func TestShouldCountSingleStringByDataCodingCorrectly(t *testing.T) {
	rawMessage := &IRawMessage{
		Recipients: "886987654321",
		Originator: "Jimmy",
		Body:       "我",
	}

	messageBuilder := &IMessageBuilder{
		IRawMessage:       rawMessage,
		ShouldConcatenate: false,
		Params: sms.Params{
			Type:       consts.SMSType,
			DataCoding: "unicode",
		},
	}

	rawCount := len(rawMessage.Body)
	expectedCount := utf8.RuneCountInString(rawMessage.Body)
	actualCount := messageBuilder.countSingleStringByDataCoding(rawCount)

	if expectedCount != actualCount {
		t.Errorf("Expected single string to be count as %d but got %d", expectedCount, actualCount)
	}
}
