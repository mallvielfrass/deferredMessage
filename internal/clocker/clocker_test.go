package clocker

import (
	"deferredMessage/internal/models"
	"fmt"
	"testing"
	"time"
)

type testService struct {
}
type poolService struct {
}

func (t poolService) GetMsgList(period time.Duration) []models.Message {
	return []models.Message{}
}
func (s testService) Send(msg models.Message) error {
	fmt.Printf("Message: %s\n", msg.Message)
	return nil
}

func TestClocker(t *testing.T) {
	tserv := testService{}
	poolService := poolService{}
	clock := NewClocker(tserv, poolService)
	clock.Start()
	clock.AddMessage(models.Message{
		Message: "hello",
		Time:    time.Now(),
	})
	time.Sleep(20 * time.Second)
	clock.Stop()
}
