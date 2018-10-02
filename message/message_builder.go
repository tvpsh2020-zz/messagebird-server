package message

import (
	"unicode/utf8"

	"github.com/messagebird/go-rest-api/sms"
	"github.com/tvpsh2020/messagebird-server/consts"
)

// IMessageBuilder is an interface for message builder
type IMessageBuilder struct {
	IRawMessage       *IRawMessage
	Recipients        []string
	Params            sms.Params
	BodyLength        int
	ShouldConcatenate bool
}

var concatenatedSMSLength = map[string]int{
	"plain":   consts.PlainCSMSLength,
	"unicode": consts.UnicodeCSMSLength,
}

var singleSMSLength = map[string]int{
	"plain":   consts.PlainSMSLength,
	"unicode": consts.UnicodeSMSLength,
}

var singleSMSMaxLength = map[string]int{
	"plain":   consts.PlainSMSMaxLength,
	"unicode": consts.UnicodeSMSMaxLength,
}

func (mb *IMessageBuilder) buildMessages() []IMessage {
	var messages []IMessage
	body := mb.IRawMessage.Body

	if mb.ShouldConcatenate {
		_buffer := ""
		_count := 0

		for len(body) > 0 {
			_, _size := utf8.DecodeRuneInString(body)
			_count += mb.countSingleStringByDataCoding(_size)

			if _count >= concatenatedSMSLength[mb.Params.DataCoding] {
				messages = append(messages, mb.makeSingleMessage(_buffer))
				_count = mb.countSingleStringByDataCoding(_size)
				_buffer = ""
			}

			_buffer += body[:_size]
			body = body[_size:]
		}

		messages = append(messages, mb.makeSingleMessage(_buffer))
	} else {
		messages = append(messages, mb.makeSingleMessage(body))
	}

	return mb.addTypeDetails(messages)
}

func (mb *IMessageBuilder) start() ([]IMessage, error) {
	if err := mb.validateOriginator(); err != nil {
		return nil, err
	}

	if err := mb.splitRecipients(); err != nil {
		return nil, err
	}

	if err := mb.checkBodyDataCoding(); err != nil {
		return nil, err
	}

	if err := mb.countBodyLength(); err != nil {
		return nil, err
	}

	return mb.buildMessages(), nil
}
