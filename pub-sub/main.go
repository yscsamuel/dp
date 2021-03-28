package main

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Pubsub struct {
	sync.RWMutex
	subscribes map[string][]chan string
}

// Constructor
func GetNewPubsub() *Pubsub {
	return &Pubsub{
		subscribes: make(map[string][]chan string),
	}
}

func (ps *Pubsub) Subscribe(channel string, buffer int) chan string {
	ps.Lock()
	defer ps.Unlock()
	// add a new subscription channel
	ch := make(chan string, buffer)
	ps.subscribes[channel] = append(ps.subscribes[channel], ch)

	return ch
}

func (ps *Pubsub) Publish(channel, msg string) {
	// guards fields subscribes read
	ps.RLock()
	defer ps.RUnlock()

	// broadcast to each channel
	for _, ch := range ps.subscribes[channel] {
		ch <- msg
		log.Infof("channel %v subscribed ", channel)

	}
}

func Listener(ch chan string) {
	for msg := range ch {
		log.Info("got message: ", msg)
	}
}

func main() {
	ps := GetNewPubsub()
	ch1 := ps.Subscribe("Hello", 1)
	ch2 := ps.Subscribe("World", 1)

	go Listener(ch1)
	go Listener(ch2)

	publish := func(ch, msg string) {
		log.Info("publish message: ", msg)
		ps.Publish(ch, msg)
		time.Sleep(1 * time.Millisecond)
	}

	time.Sleep(1 * time.Millisecond)
	publish("Hello", "hello message")
	publish("World", "world message")
	publish("World", "world message2")
	time.Sleep(1 * time.Millisecond)
}
