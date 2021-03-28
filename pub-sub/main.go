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

const buf int = 1

// Constructor
func GetNewPubsub() *Pubsub {
	return &Pubsub{
		subscribes: make(map[string][]chan string),
	}
}

func (ps *Pubsub) Subscribe(channel string) chan string {
	ps.Lock()
	defer ps.Unlock()
	// add a new subscription channel
	ch := make(chan string, buf)
	log.Infof("channel %v subscribed ", channel)
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
	}
}

func Listener(ch <-chan string) {
	// blocking until receive message
	log.Info("got message: ", <-ch)
}

func main() {
	ps := GetNewPubsub()
	ch1 := ps.Subscribe("Hello")
	ch2 := ps.Subscribe("World")

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
