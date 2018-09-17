package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/messagebird/go-rest-api/sms"
	"github.com/tvpsh2020/messagebird-server/consts"
)

// IMessageBuilder is an interface for message builder
type IMessageBuilder struct {
	IRawMessage   *IRawMessage
	Recipients    []string
	Params        sms.Params
	BodyLength    int
	SplitParts    int
	SplitPartSize int
}

var concatenatedSMSLength = map[string]int{
	"plain":   consts.PlainCSMSLength,
	"unicode": consts.UnicodeCSMSLength,
}

var singleSMSLength = map[string]int{
	"plain":   consts.PlainSMSLength,
	"unicode": consts.UnicodeSMSLength,
}

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

	if len(result) == 0 {
		return errors.New("recipients is illegal")
	}

	mb.Recipients = result

	return nil
}

func (mb *IMessageBuilder) checkDataCoding() {
	var regex = regexp.MustCompile(consts.GSM0338Regex)

	for _, match := range regex.FindAllString(mb.IRawMessage.Body, -1) {
		if len(match) > 0 {
			mb.Params.DataCoding = consts.Unicode
			log.Println("DataCoding: ", mb.Params.DataCoding)
			return
		}
	}

	mb.Params.DataCoding = consts.Plain

	log.Println("DataCoding: ", mb.Params.DataCoding)
}

func (mb *IMessageBuilder) countBody() error {

	mb.IRawMessage.Body = strings.TrimSpace(mb.IRawMessage.Body)

	mb.BodyLength = len(mb.IRawMessage.Body)

	if mb.BodyLength == 0 {
		return errors.New("message content is illegal")
	}

	mb.checkDataCoding()

	if mb.BodyLength > singleSMSLength[mb.Params.DataCoding] {
		mb.SplitPartSize = concatenatedSMSLength[mb.Params.DataCoding]
		mb.SplitParts = mb.BodyLength / mb.SplitPartSize

		if mb.BodyLength%mb.SplitPartSize > 0 {
			mb.SplitParts++
		}
	} else {
		mb.SplitPartSize = singleSMSLength[mb.Params.DataCoding]
		mb.SplitParts = 1
	}

	fmt.Println("Count -> ", mb.BodyLength)
	return nil
}

func (mb *IMessageBuilder) stringToBinary(str string) string {
	src := []byte(str)
	encodedStr := hex.EncodeToString(src)

	return encodedStr
}

func (mb *IMessageBuilder) buildMessages() []IMessage {
	var result []IMessage
	rand.Seed(time.Now().UTC().UnixNano())
	referenceNum := rand.Intn(256)

	for i := 0; i < mb.SplitParts; i++ {
		typeDetails := make(sms.TypeDetails)
		typeDetails["udh"] = fmt.Sprintf("050003%02x%02x%02x", referenceNum, mb.SplitParts, i+1)

		mb.Params.TypeDetails = typeDetails

		_body := ""

		if mb.SplitParts == 1 {
			_body = mb.IRawMessage.Body
		} else {
			if len(mb.IRawMessage.Body[i*mb.SplitPartSize:]) < mb.SplitPartSize {
				_body = mb.IRawMessage.Body[i*mb.SplitPartSize:]
			} else {
				_body = mb.IRawMessage.Body[i*mb.SplitPartSize : (i+1)*mb.SplitPartSize]
			}
		}

		_result := &IMessage{
			Originator: mb.IRawMessage.Originator,
			// Body:       mb.stringToBinary(_body),
			Body:       _body,
			Recipients: mb.Recipients,
			Params:     mb.Params,
		}

		result = append(result, *_result)
	}

	return result
}

func (mb *IMessageBuilder) start() ([]IMessage, error) {
	if err := mb.validateOriginator(); err != nil {
		return nil, err
	}

	if err := mb.splitRecipients(); err != nil {
		return nil, err
	}

	if err := mb.countBody(); err != nil {
		return nil, err
	}

	return mb.buildMessages(), nil
}
