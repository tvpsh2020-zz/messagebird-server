package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/messagebird/go-rest-api/sms"
)

type MessageBuilder struct {
	RawMessageBody *RawMessageBody
	Recipients     []string
	Params         sms.Params
	BodyLength     int
	SplitParts     int
	SplitPartSize  int
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

func (mb *MessageBuilder) constructor() {
	mb.Params = sms.Params{
		Type: "binary",
	}
}

func (mb *MessageBuilder) splitRecipients() {
	removeWhiteSpace := strings.Replace(mb.RawMessageBody.Recipients, " ", "", -1)
	mb.Recipients = strings.Split(removeWhiteSpace, ",")
}

func (mb *MessageBuilder) countBody() {
	mb.BodyLength = len(mb.RawMessageBody.Body)

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
			_body = mb.RawMessageBody.Body
		} else {
			if len(mb.RawMessageBody.Body[i*mb.SplitPartSize:]) < mb.SplitPartSize {
				_body = mb.RawMessageBody.Body[i*mb.SplitPartSize:]
			} else {
				_body = mb.RawMessageBody.Body[i*mb.SplitPartSize : (i+1)*mb.SplitPartSize]
			}
		}

		_result := &Message{
			Originator: mb.RawMessageBody.Originator,
			Body:       _body,
			Recipients: mb.Recipients,
			Params:     mb.Params,
		}

		result = append(result, *_result)
	}

	return result
}

func (mb *MessageBuilder) start() []Message {
	mb.constructor()
	mb.splitRecipients()
	mb.countBody()

	return mb.buildMessages()
}
