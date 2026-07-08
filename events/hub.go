package events

import "sync"

type Hub struct {
	mutex       sync.RWMutex
	subscribers map[string]map[*subscriber]struct{}
}

func NewHub() *Hub {
	return &Hub{subscribers: make(map[string]map[*subscriber]struct{})}
}

var Default = NewHub()

func (self *Hub) Subscribe(topic string) (<-chan Event, func()) {
	member := &subscriber{channel: make(chan Event, SubscriberBuffer)}

	self.mutex.Lock()
	if self.subscribers[topic] == nil {
		self.subscribers[topic] = make(map[*subscriber]struct{})
	}
	self.subscribers[topic][member] = struct{}{}
	self.mutex.Unlock()

	unsubscribe := func() {
		self.mutex.Lock()
		defer self.mutex.Unlock()

		set, present := self.subscribers[topic]
		if !present {
			return
		}
		if _, exists := set[member]; !exists {
			return
		}

		delete(set, member)
		close(member.channel)
		if len(set) == 0 {
			delete(self.subscribers, topic)
		}
	}

	return member.channel, unsubscribe
}

func (self *Hub) Publish(topic string, event Event) {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	for member := range self.subscribers[topic] {
		select {
		case member.channel <- event:
		default:
		}
	}
}
