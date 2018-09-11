package main

import (
	"log"
	"time"
)

func initTaskManager() {
	checkMessageQueue := time.Tick(time.Second * 1)

	go func() {
		for {
			select {
			case <-checkMessageQueue:
				MessageQueue.RLock()

				if len(MessageQueue.List) > 0 {
					log.Println("Sent message.")
					log.Println("MQ List Count : ", len(MessageQueue.List))
					SendSMSToMessageBirdV2(MessageQueue.List[0])
					MessageQueue.RUnlock()
					MessageQueue.Lock()
					MessageQueue.List = MessageQueue.List[1:]
					MessageQueue.Unlock()
				} else {
					MessageQueue.RUnlock()
				}
			}
		}
	}()
}
