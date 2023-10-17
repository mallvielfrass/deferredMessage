package clocker

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/ltcsuite/lnd/ticker"
)

type Sender interface {
	Send(msg Message) error
}
type Clocker struct {
	messages            []Message
	messagesBuffer      chan Message
	deleteMessageBuffer chan Message
	sender              Sender
	ticker              *ticker.T
	done                chan bool
	mutex               sync.Mutex
}

func NewClocker(sender Sender) *Clocker {
	return &Clocker{
		messages:            []Message{},
		messagesBuffer:      make(chan Message, 100),
		deleteMessageBuffer: make(chan Message, 100),
		sender:              sender,
		done:                make(chan bool),
	}
}

// Start
func (c *Clocker) Start() {
	fmt.Println("Clocker started")
	c.ticker = ticker.New(100 * time.Millisecond)
	go c.tick()
}
func (c *Clocker) handleMessages(messages []Message) {
	for _, msg := range messages {
		go c.sender.Send(msg)
	}
}
func (c *Clocker) checkMsgTime() {
	c.ticker.Pause()
	c.mutex.Lock()
	fmt.Printf("Messages: %v\n", c.messages)
	messages, index := binarySearchMsgWithTime(c.messages, time.Now())
	if 0 < len(messages) {
		c.handleMessages(messages)
		c.messages = c.messages[index:]
	}
	c.mutex.Unlock()
	c.ticker.Resume()
}
func (c *Clocker) tick() {
	c.ticker.Resume()
	for {
		select {
		case val := <-c.messagesBuffer:
			fmt.Printf("New Message: %s\n", val.Message)
			c.addMessage(val)
		case val := <-c.deleteMessageBuffer:
			fmt.Printf("Delete Message: %s\n", val.Message)
			c.deleteMessage(val)
		case <-c.done:
			fmt.Println("Clocker stopped")
			return
		case t := <-c.ticker.Ticks():
			fmt.Println("Tick at", t)
			c.checkMsgTime()

		}
	}
}

// Stop
func (c *Clocker) Stop() {
	fmt.Println("Stopping clocker")
	c.ticker.Stop()
	c.done <- true
	fmt.Println("Clocker stopped")
}
func (c *Clocker) addMessage(msg Message) {
	c.ticker.Pause()
	c.mutex.Lock()
	c.messages = append(c.messages, msg)
	sort.Slice(c.messages, func(i, j int) bool {
		return c.messages[i].Time.Before(c.messages[j].Time)
	})
	c.mutex.Unlock()
	c.ticker.Resume()
}

// delete message
func (c *Clocker) deleteMessage(msg Message) {
	c.ticker.Pause()
	c.mutex.Lock()
	c.messages = findAndRemoveMessage(c.messages, msg)
	c.mutex.Unlock()
	c.ticker.Resume()
}

// External func
// add message
func (c *Clocker) AddMessage(msg Message) {
	c.messagesBuffer <- msg
}

// Delete message
func (c *Clocker) DeleteMessage(msg Message) {
	c.deleteMessageBuffer <- msg
}
