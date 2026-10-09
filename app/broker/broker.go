// Package broker ...
package broker

import (
	"sync"

	"my-redis/app/client"
	"my-redis/app/parser"
)

type Broker struct {
	mu       sync.Mutex
	channels map[string]map[*client.Client]struct{}
	byClient map[*client.Client]map[string]struct{}
}

func New() *Broker {
	return &Broker{
		channels: map[string]map[*client.Client]struct{}{},
		byClient: map[*client.Client]map[string]struct{}{},
	}
}

func (b *Broker) Subscribe(c *client.Client, ch string) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	subs, ok := b.byClient[c]
	if !ok {
		subs = map[string]struct{}{}
		b.byClient[c] = subs
	}
	subs[ch] = struct{}{}
	cls, ok := b.channels[ch]
	if !ok {
		cls = map[*client.Client]struct{}{}
		b.channels[ch] = cls
	}
	cls[c] = struct{}{}
	return len(subs)
}

func (b *Broker) Unsubscribe(c *client.Client, chans ...string) int {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(chans) == 0 { // unsubscribe from all
		for ch := range b.byClient[c] {
			chans = append(chans, ch)
		}
		if len(chans) == 0 {
			return 0
		}

	}
	for _, ch := range chans {
		b.unsubscribeLocked(c, ch)
	}
	return len(b.byClient[c])
}

func (b *Broker) unsubscribeLocked(c *client.Client, ch string) {
	subs, ok := b.byClient[c]
	if !ok {
		return
	}
	delete(subs, ch)
	if len(subs) == 0 {
		delete(b.byClient, c)
	}
	clients := b.channels[ch]
	delete(clients, c)
	if len(clients) == 0 {
		delete(b.channels, ch)
	}
}

func (b *Broker) Publish(ch, msg string) int {
	data := parser.StringsToBytes("message", ch, msg)
	b.mu.Lock()
	defer b.mu.Unlock()
	var n int
	for c := range b.channels[ch] {
		if c.Reply(data) {
			n++
		}
	}
	return n
}
