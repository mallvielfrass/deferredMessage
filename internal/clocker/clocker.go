package clocker

import (
	"deferredMessage/internal/models"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/ltcsuite/lnd/ticker"
)

type Sender interface {
	Send(msg models.Message) error
}
type Pool interface {
	GetMsgList(period time.Duration) []models.Message
}
type Clocker struct {
	messages            []models.Message
	messagesBuffer      chan models.Message
	deleteMessageBuffer chan models.Message
	sender              Sender
	pool                Pool
	ticker              *ticker.T
	poolTicker          *ticker.T
	done                chan bool
	mutex               sync.Mutex
}

func NewClocker(sender Sender, pool Pool) *Clocker {
	return &Clocker{
		messages:            []models.Message{},
		messagesBuffer:      make(chan models.Message, 100),
		deleteMessageBuffer: make(chan models.Message, 100),
		sender:              sender,
		pool:                pool,
		done:                make(chan bool),
	}
}

// Start
func (c *Clocker) Start() {
	fmt.Println("Clocker started")
	c.ticker = ticker.New(100 * time.Millisecond)
	c.poolTicker = ticker.New(10 * time.Minute)
	go c.tick()
}
func (c *Clocker) handleMessages(messages []models.Message) {
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

// synsPool
func (c *Clocker) synsPool() {
	c.ticker.Pause()
	c.mutex.Lock()
	newMsgPool := c.pool.GetMsgList(10 * time.Minute)
	sort.Slice(newMsgPool, func(i, j int) bool {
		return newMsgPool[i].Time.Before(newMsgPool[j].Time)
	})
	c.messages = newMsgPool
	c.mutex.Unlock()
	c.ticker.Resume()
}
func (c *Clocker) tick() {
	c.ticker.Resume()
	c.poolTicker.Resume()
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
		case <-c.poolTicker.Ticks():
			fmt.Println("pool Tick")
			if len(c.messages) > 0 {
				c.synsPool()
			}
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
	c.poolTicker.Stop()
	c.done <- true
	fmt.Println("Clocker stopped")
}
func (c *Clocker) addMessage(msg models.Message) {
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
func (c *Clocker) deleteMessage(msg models.Message) {
	c.ticker.Pause()
	c.mutex.Lock()
	c.messages = findAndRemoveMessage(c.messages, msg)
	c.mutex.Unlock()
	c.ticker.Resume()
}

// External func
// add message
func (c *Clocker) AddMessage(msg models.Message) {
	c.messagesBuffer <- msg
}

// Delete message
func (c *Clocker) DeleteMessage(msg models.Message) {
	c.deleteMessageBuffer <- msg
}
