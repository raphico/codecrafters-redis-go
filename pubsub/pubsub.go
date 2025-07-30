package pubsub

import "sync"

type subscriber interface {
	SendMessage(channel, message string)
	GetSessionID() string
}

type subscriptions map[string]map[subscriber]struct{}
type subsByClient map[subscriber]map[string]struct{}

type PubsubManager struct {
	// to prevent crashes and race conditions when multiple clients subscribe to the same channel concurrently
	mu sync.RWMutex

	// we are keeping track of two mappings for faster publish, unsubscribe, subscription count
	// tradeoff: more efficient operation, but uses slightly more memory
	subscriptions subscriptions
	subsByClient  subsByClient
}

func NewPubSubManager() *PubsubManager {
	return &PubsubManager{
		subscriptions: make(subscriptions),
		subsByClient:  make(subsByClient),
	}
}

func (ps *PubsubManager) Subscribe(channel string, sub subscriber) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, ok := ps.subscriptions[channel]; !ok {
		ps.subscriptions[channel] = make(map[subscriber]struct{})
	}

	if _, ok := ps.subsByClient[sub]; !ok {
		ps.subsByClient[sub] = make(map[string]struct{})
	}

	ps.subscriptions[channel][sub] = struct{}{}
	ps.subsByClient[sub][channel] = struct{}{}
}

func (ps *PubsubManager) GetSubscribedCount(sub subscriber) int {
	// reading or writing to the same map concurrently is not safe(race conditions, panics) in go
	// so we are blocking all readers and writers while getting the subscriber count
	ps.mu.Lock()
	defer ps.mu.Unlock()

	channels, ok := ps.subsByClient[sub]
	if !ok {
		return 0
	}

	return len(channels)
}

func (ps *PubsubManager) Publish(channel, message string) int {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	subs, ok := ps.subscriptions[channel]
	if !ok {
		return 0
	}

	for sub := range subs {
		sub.SendMessage(channel, message)
	}

	return len(subs)
}

func (ps *PubsubManager) Unsubscribe(channel string, sub subscriber) int {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	removed := 0

	if subs, ok := ps.subscriptions[channel]; ok {
		if _, subscribed := subs[sub]; subscribed {
			delete(subs, sub)
			removed = 1
			if len(subs) == 0 {
				delete(ps.subscriptions, channel)
			}
		}
	}

	if channels, ok := ps.subsByClient[sub]; ok {
		if _, exists := channels[channel]; exists {
			delete(channels, channel)
			if len(channels) == 0 {
				delete(ps.subsByClient, sub)
			}
		}
	}

	return removed
}
