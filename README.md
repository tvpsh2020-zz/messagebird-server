# Jimmmmy SMS Server

Sending SMS by integrate MessageBird service. The progress is about 80% now.

## Requirements

- Sign up for a free MessageBird account
- Top-up amount of credit you'd like to add to your balance
- Create a new access key in the developers sections, then add key to the config file

## Installation

```
$ go get github.com/tvpsh2020/messagebird-server
```

## Usage

### 1. Send SMS

`POST /api/message`

with body
```
- originator: Sender identity, this can be a telephone number (including country code) or an alphanumeric string. In case of an alphanumeric string, the maximum length is 11 characters.
- body: The body of the SMS message, should not be empty.
- recipients: Receivers, multiple receiver should be split with comma, any words other than number will be filter out.
```

body example
```
{
	"originator": "12345678901",
	"body": "\"abcd",
	"recipients": "886987654321, +886987654321"
}
```

## Developer's Note

### Noun define

- server: mean this API server

### Rule of MessageBird

- If message content/body is longer than 160 chars with plain format (70 with unicode format), split it into multiple parts. The rule of split mechanism is define in [here](https://support.messagebird.com/hc/en-us/articles/208739745-How-long-is-1-SMS-Message).
- Do not send empty or incorrect parameter values to MessageBird
- Throughput of MessageBird API is one API request per second, no matter what jobs you are doing
- Sending SMS with [UDH](https://en.wikipedia.org/wiki/Concatenated_SMS) (User Data Header) can make MessageBird Server know that your SMS should be concatenated.
- In theory the concatenated SMS may consist of up to 255 separate SMS messages, but in MessageBird, the concatenated SMS is maximum separated to 9 parts, if your SMS over the rule, the whole messages will be buffered on MessageBird Server and THIS WILL BE CHARGED WITHOUT SENDING SUCCESSFUL except you set the validity time.
- Euro sign is not a Plain text for MessageBird, should be send with Unicode to display.
- Curly Quotation Marks in Plain text mode will be translate, should be send with Unicode to display.
```
‘’ -> '
“” -> "
```


### Done work

- Make config loader to prevent api-key exposed
- Get balance from MessageBird
- Send test message for testing connection with MessageBird
- Send message from server
- MessageBird can only accept 1 req per second, let all received messages be stored in a queue
- Make a http route to handle flexible uri pattern
- Make a task management to handle whole server event
- Make config loader easy to load and easy to add new parameter
- Update MessageBird API from v4 to v5
- Validate message originator
- Validate recipients
- Message body check: message body contain unicode or only plain
- Message body check: > 160 (or unicode > 70) should be split
- Message body check: should not be empty
- Message body check: invalid word

### Todo
- API validate
- Uniform API log
- Memory usage
- Test
- Refactor into OOP

### Nice to have
- Add validity time to SMS
- Add reportUrl to receive status of sent SMS

### Knowing issue
- Can only send hex string with plain text, if you send with unicode, it will be fail.
