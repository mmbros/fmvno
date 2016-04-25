package fw

import (
	"sync"

	"github.com/mmbros/fmvno/model"
	"github.com/mmbros/fmvno/util"
)

// QueueSpedizioni represents a thread safe queue of spedizioni
type QueueSpedizioni struct {
	buf util.Queue
	mx  sync.RWMutex
}

// NewQueueSpedizioni return a new QueueSpedizioni with given initial capacity
func NewQueueSpedizioni(initCap int) *QueueSpedizioni {
	q := &QueueSpedizioni{
		buf: util.NewCircularFifoQueue(initCap),
		mx:  sync.RWMutex{},
	}
	return q
}

// Push ...
func (q *QueueSpedizioni) Push(item *model.Spedizione) {
	q.mx.Lock()
	defer q.mx.Unlock()
	q.buf.Push(item)
}

// Pop ...
func (q *QueueSpedizioni) Pop() *model.Spedizione {
	q.mx.Lock()
	defer q.mx.Unlock()
	return q.buf.Pop().(*model.Spedizione)
}

// Len ...
func (q *QueueSpedizioni) Len() int {
	q.mx.RLock()
	defer q.mx.RUnlock()
	return q.buf.Len()
}

// Cap ...
func (q *QueueSpedizioni) Cap() int {
	q.mx.RLock()
	defer q.mx.RUnlock()
	return q.buf.Cap()
}

// HandleReceive ...
func (q *QueueSpedizioni) HandleReceive(input <-chan *model.Spedizione) {
	for sped := range input {
		q.Push(sped)
	}
}

// HandleSend ...
func (q *QueueSpedizioni) HandleSend(output chan<- *model.Spedizione, maxSped int) {
	L := q.Len()
	if L > maxSped {
		L = maxSped
	}
	for j := 0; j < L; j++ {
		output <- q.Pop()
	}
}
