# Jimmmmy SMS Server

Sending SMS by integrate MessageBird service. The progress is about 50% now.

## Requirements

- Sign up for a free MessageBird account
- Top-up amount of credit you'd like to add to your balance
- Create a new access key in the developers sections, then add key to the config file

## Installation

```
$ go get github.com/tvpsh2020/messagebird-server
```

## Developer's Note

### Noun define

- server: mean this API server

### Rule from MessageBird

- If message content/body is longer than 160 chars with plain format (70 with unicode format), split it into multiple parts. The rule of split mechanism is define in [here](https://support.messagebird.com/hc/en-us/articles/208739745-How-long-is-1-SMS-Message)
- Do not send empty or incorrect parameter values to MessageBird
- Throughput of MessageBird API is one API request per second, no matter what jobs you are doing
- Sending SMS with [UDH](https://en.wikipedia.org/wiki/Concatenated_SMS) (User Data Header) can make MessageBird Server know that your SMS should be concatenated.
- In theory the concatenated SMS may consist of up to 255 separate SMS messages, but in MessageBird, the concatenated SMS is maximum separated to 9 parts, if your SMS over the rule, the whole messages will be buffered on MessageBird Server and THIS WILL BE CHARGED WITHOUT SENDING SUCCESSFUL.

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

### Todo

- Message body check: message body contain unicode or only plain
- Message body check: > 160 (or unicode > 70) should be split
- Message body check: should not be empty
- Message body check: invalid word
- (Last step) Refactor into OOP
