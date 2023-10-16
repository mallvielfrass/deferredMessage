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
	messages       []Message
	messagesBuffer chan Message
	sender         Sender
	ticker         *ticker.T
	done           chan bool
	mutex          sync.Mutex
}

func NewClocker(sender Sender) *Clocker {
	return &Clocker{
		messages:       []Message{},
		messagesBuffer: make(chan Message),
		sender:         sender,
		done:           make(chan bool),
	}
}

// Start
func (c *Clocker) Start() {
	fmt.Println("Clocker started")
	c.ticker = ticker.New(500 * time.Millisecond)
	go c.tick()
}

func (c *Clocker) checkMsgTIme() {
	c.ticker.Pause()
	c.mutex.Lock()
	fmt.Printf("Messages: %v\n", c.messages)
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
		case <-c.done:
			fmt.Println("Clocker stopped")
			return
		case t := <-c.ticker.Ticks():
			fmt.Println("Tick at", t)
			c.checkMsgTIme()

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

// add message
func (c *Clocker) AddMessage(msg Message) {
	c.messagesBuffer <- msg
}
