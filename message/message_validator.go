package message

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/messagebird/go-rest-api/sms"
	"github.com/tvpsh2020/messagebird-server/consts"
)

func (mb *IMessageBuilder) validateOriginator() error {
	if !regexp.MustCompile(consts.OriginatorRegex).MatchString(mb.IRawMessage.Originator) {
		return errors.New("originator is illegal")
	}

	return nil
}

func (mb *IMessageBuilder) fixRecipients() string {
	var regex = regexp.MustCompile(consts.RecipientsRegex)
	var validatedString string

	for _, match := range regex.FindAllString(mb.IRawMessage.Recipients, -1) {
		if len(match) > 0 {
			validatedString += match
		}
	}

	fixedRecipients := strings.Replace(validatedString, " ", "", -1)
	return fixedRecipients
}

func (mb *IMessageBuilder) splitRecipients() error {
	fixedRecipients := mb.fixRecipients()
	splitWithComma := strings.Split(fixedRecipients, ",")

	var result []string

	for _, str := range splitWithComma {
		if str != "" {
			result = append(result, str)
		}
	}

	if len(result) == 0 || len(result) > 50 {
		return errors.New("recipients is illegal")
	}

	mb.Recipients = result

	return nil
}

func (mb *IMessageBuilder) checkBodyDataCoding() error {
	mb.IRawMessage.Body = strings.TrimSpace(mb.IRawMessage.Body)

	if len(mb.IRawMessage.Body) == 0 {
		return errors.New("message content cannot be empty")
	}

	var regex = regexp.MustCompile(consts.GSM0338Regex)

	for _, match := range regex.FindAllString(mb.IRawMessage.Body, -1) {
		if len(match) > 0 {
			mb.Params.DataCoding = consts.Unicode
			return nil
		}
	}

	mb.Params.DataCoding = consts.Plain

	if len(mb.IRawMessage.Body) > singleSMSMaxLength[mb.Params.DataCoding] {
		return errors.New("message content is too long")
	}

	return nil
}

func (mb *IMessageBuilder) countBodyLength() error {
	switch mb.Params.DataCoding {
	case consts.Plain:
		mb.BodyLength = len(mb.IRawMessage.Body)
	case consts.Unicode:
		mb.BodyLength = utf8.RuneCountInString(mb.IRawMessage.Body)
	default:
		return errors.New("undefined data coding")
	}

	if mb.BodyLength > singleSMSLength[mb.Params.DataCoding] {
		mb.ShouldConcatenate = true
	}

	return nil
}

func (mb *IMessageBuilder) makeSingleMessage(body string) IMessage {
	result := IMessage{
		Originator: mb.IRawMessage.Originator,
		Body:       body,
		Recipients: mb.Recipients,
		Params:     mb.Params,
	}

	return result
}

func (mb *IMessageBuilder) countSingleStringByDataCoding(rawSize int) int {
	if mb.Params.DataCoding == consts.Unicode {
		return 1
	}

	return rawSize
}

func (mb *IMessageBuilder) addTypeDetails(messages []IMessage) []IMessage {
	rand.Seed(time.Now().UTC().UnixNano())
	referenceNum := rand.Intn(256)
	var resultMessages []IMessage

	for i, message := range messages {
		typeDetails := make(sms.TypeDetails)
		typeDetails["udh"] = fmt.Sprintf("050003%02x%02x%02x", referenceNum, len(messages), i+1)

		message.Params.TypeDetails = typeDetails
		resultMessages = append(resultMessages, message)
	}

	return resultMessages
}
