package pubsub

import "sync"

type Agent[T any] struct {
	mu          sync.Mutex
	subscribers map[Topic][]chan T
	quit        chan struct{}
	closed      bool
}

func NewAgent[T any]() *Agent[T] {
	return &Agent[T]{
		subscribers: make(map[Topic][]chan T),
		quit:        make(chan struct{}),
	}
}

func (a *Agent[T]) Publish(topic Topic, data T) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.closed {
		return
	}

	for _, subscriber := range a.subscribers[topic] {
		subscriber <- data
	}
}

func (a *Agent[T]) Subscribe(topic Topic) <-chan T {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.closed {
		return nil
	}

	if _, ok := a.subscribers[topic]; !ok {
		a.subscribers[topic] = make([]chan T, 0)
	}

	ch := make(chan T, 1024)
	a.subscribers[topic] = append(a.subscribers[topic], ch)

	return ch
}

func (b *Agent[T]) Close() {
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
