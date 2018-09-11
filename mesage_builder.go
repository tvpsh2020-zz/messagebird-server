package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/messagebird/go-rest-api/sms"
)

type MessageBuilder struct {
	RawMessage    *RawMessage
	Recipients    []string
	Params        sms.Params
	BodyLength    int
	SplitParts    int
	SplitPartSize int
}

const (
	PlainSMSLength    = 160
	PlainCSMSLength   = 153
	UnicodeSMSLength  = 70
	UnicodeCSMSLength = 67
)

var concatenatedSMSLength = map[string]int{
	"plain":   PlainCSMSLength,
	"unicode": UnicodeCSMSLength,
}

var singleSMSLength = map[string]int{
	"plain":   PlainSMSLength,
	"unicode": UnicodeSMSLength,
}

func (mb *MessageBuilder) validateOriginator() error {
	if !regexp.MustCompile(`^[a-zA-Z0-9]{1,11}$`).MatchString(mb.RawMessage.Originator) {
		return errors.New("originator is illegal")
	}

	return nil
}

func (mb *MessageBuilder) splitRecipients() error {
	// var rule = regexp.MustCompile(`[0-9]*\,*`)
	// var result string

	// for _, match := range re.FindAllString(str, -1) {
	// 	if len(match) > 0 {
	// 		fmt.Println(match, "found at index", i)
	// 		result += match
	// 	}

	// }

	removeWhiteSpace := strings.Replace(mb.RawMessage.Recipients, " ", "", -1)
	mb.Recipients = strings.Split(removeWhiteSpace, ",")
	return nil
}

func (mb *MessageBuilder) countBody() {
	mb.BodyLength = len(mb.RawMessage.Body)

	// check empty body

	// if is plain text, set here
	mb.Params.DataCoding = "plain"

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
}

func (mb *MessageBuilder) stringToBinary(str string) string {
	src := []byte(str)
	encodedStr := hex.EncodeToString(src)

	return encodedStr
}

func (mb *MessageBuilder) buildMessages() []Message {
	var result []Message
	rand.Seed(time.Now().UTC().UnixNano())
	referenceNum := rand.Intn(256)

	for i := 0; i < mb.SplitParts; i++ {
		typeDetails := make(sms.TypeDetails)
		typeDetails["udh"] = fmt.Sprintf("050003%02x%02x%02x", referenceNum, mb.SplitParts, i+1)

		mb.Params.TypeDetails = typeDetails

		_body := ""

		if mb.SplitParts == 1 {
			_body = mb.RawMessage.Body
		} else {
			if len(mb.RawMessage.Body[i*mb.SplitPartSize:]) < mb.SplitPartSize {
				_body = mb.RawMessage.Body[i*mb.SplitPartSize:]
			} else {
				_body = mb.RawMessage.Body[i*mb.SplitPartSize : (i+1)*mb.SplitPartSize]
			}
		}

		_result := &Message{
			Originator: mb.RawMessage.Originator,
			Body:       mb.stringToBinary(_body),
			Recipients: mb.Recipients,
			Params:     mb.Params,
		}

		result = append(result, *_result)
	}

	return result
}

func (mb *MessageBuilder) start() ([]Message, error) {
	if err := mb.validateOriginator(); err != nil {
		return nil, err
	}

	mb.splitRecipients()
	mb.countBody()

	return mb.buildMessages(), nil
}
