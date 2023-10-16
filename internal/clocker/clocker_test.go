package clocker

import (
	"fmt"
	"testing"
	"time"
)

type testService struct {
}

func (s testService) Send(msg Message) error {
	fmt.Printf("Message: %s\n", msg.Message)
	return nil
}

func TestClocker(t *testing.T) {
	tserv := testService{}
	clock := NewClocker(tserv)
	clock.Start()
	clock.AddMessage(Message{
		Message: "hello",
		Time:    time.Now(),
	})
	time.Sleep(20 * time.Second)
	clock.Stop()
}
