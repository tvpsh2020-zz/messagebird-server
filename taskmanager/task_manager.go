package taskmanager

import (
	"time"

	"github.com/tvpsh2020/messagebird-server/message"
)

// Initialize task manager
func Initialize() {
	checkMessageQueue := time.Tick(time.Second * 1)

	go func() {
		for {
			select {
			case <-checkMessageQueue:
				message.SendFromMessageQueue()
			}
		}
	}()
}
