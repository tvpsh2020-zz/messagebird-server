# Jimmmmy SMS Server

This is a 3rd party server for sending SMS by integrate MessageBird service. The progress is about 50% now.

# Developer memo

## Noun define

- server: mean this API server

## Rule from MessageBird

- If message content/body is longer than 160 chars with plain format (70 with unicode format), split it into multiple parts. The rule of split mechanism is define in [here](https://support.messagebird.com/hc/en-us/articles/208739745-How-long-is-1-SMS-Message)
- Do not send empty or incorrect parameter values to MessageBird
- Throughput of MessageBird API is one API request per second, no matter what jobs you are doing
- [UDH](https://en.wikipedia.org/wiki/Concatenated_SMS) is important.

## Done work

- Make config loader to prevent api-key exposed
- Get balance from MessageBird
- Send test message for testing connection with MessageBird
- Send message from server
- MessageBird can only accept 1 req per second, let all received messages be stored in a queue
- Make a http route to handle flexible uri pattern
- Make a task management to handle whole server event
- Make config loader easy to load and easy to add new parameter

## Todo

- Message body check: message body contain unicode or only plain
- Message body check: > 160 (or unicode > 70) should be split
- Message body check: should not be empty
- Message body check: invalid word
- (Last step) Refactor into OOP
