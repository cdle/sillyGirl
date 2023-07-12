package core

import (
	"errors"
	"sync"
)

var queues sync.Map

type Queue struct {
	data []*QMessage
	head int
	tail int
	size int
	lock *sync.Mutex
}

func NewQueue(name string, size int) *Queue {
	q := &Queue{
		head: 0,
		tail: 0,
		size: size,
	}
	v, ok := queues.LoadOrStore(name, q)
	if ok {
		q = v.(*Queue)
	} else {
		q.data = make([]*QMessage, size)
		q.lock = new(sync.Mutex)
	}
	return q
}

func (q *Queue) Enqueue(value *QMessage) error {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.IsFull() {
		return errors.New("queue is full")
	}
	q.data[q.tail] = value
	q.tail = (q.tail + 1) % q.size
	return nil
}

func (q *Queue) Dequeue() (*QMessage, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.IsEmpty() {
		return nil, errors.New("queue is empty")
	}
	value := q.data[q.head]
	q.head = (q.head + 1) % q.size

	return value, nil
}

func (q *Queue) IsEmpty() bool {
	return q.head == q.tail
}

func (q *Queue) IsFull() bool {
	return (q.tail+1)%q.size == q.head
}

func (q *Queue) Size() int {
	return (q.tail - q.head + q.size) % q.size
}

func (q *Queue) GetValues() []*QMessage {
	q.lock.Lock()
	defer q.lock.Unlock()

	values := make([]*QMessage, q.Size())
	for i := 0; i < q.Size(); i++ {
		values[i] = q.data[(q.head+i)%q.size]
	}

	return values
}
