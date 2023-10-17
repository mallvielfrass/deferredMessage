package clocker

import (
	"deferredMessage/internal/models"
	"fmt"
	"testing"
	"time"
)

type testService struct {
}

func (s testService) Send(msg models.Message) error {
	fmt.Printf("Message: %s\n", msg.Message)
	return nil
}

func TestClocker(t *testing.T) {
	tserv := testService{}
	clock := NewClocker(tserv)
	clock.Start()
	clock.AddMessage(models.Message{
		Message: "hello",
		Time:    time.Now(),
	})
	time.Sleep(20 * time.Second)
	clock.Stop()
}
