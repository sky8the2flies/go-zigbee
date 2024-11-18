package pubsub

import "sync"

type Agent struct {
	mu          sync.Mutex
	subscribers map[Topic][]chan any
	quit        chan struct{}
	closed      bool
}

func NewAgent() *Agent {
	return &Agent{
		subscribers: make(map[Topic][]chan any),
		quit:        make(chan struct{}),
	}
}

func (a *Agent) Publish(topic Topic, data any) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.closed {
		return
	}

	for _, subscriber := range a.subscribers[topic] {
		subscriber <- data
	}
}

func (a *Agent) Subscribe(topic Topic) <-chan any {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.closed {
		return nil
	}

	if _, ok := a.subscribers[topic]; !ok {
		a.subscribers[topic] = make([]chan any, 0)
	}

	ch := make(chan any, 1024)
	a.subscribers[topic] = append(a.subscribers[topic], ch)

	return ch
}

func (b *Agent) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return
	}

	b.closed = true
	close(b.quit)

	for _, ch := range b.subscribers {
		for _, sub := range ch {
			close(sub)
		}
	}
}
