package util

import "github.com/mmbros/demo/queue/common"

// InfBufChan implements an infinite buffered channel
type InfBufChan struct {
	buf      Queue
	request  chan operation
	response chan interface{}
}

type operation int

const (
	opGetItem operation = iota
	opQuitNext
	opQuitDequeue
)

// NewInfBufChan returns a new InfBufChan with given initial capacity
func NewInfBufChan(initCap int) *InfBufChan {
	return &InfBufChan{
		buf:      NewCircularFifoQueue(initCap),
		request:  make(chan operation),
		response: make(chan interface{}),
	}
}

//
func (q *InfBufChan) receive(input <-chan interface{}) {
	// il chan RESP deve essere chiuso SOLO durante la gestione di una REQ

	var inputClosed, closingResp bool
	var myInput = input

	//loop:
	for {
		select {
		case n, ok := <-myInput:
			if ok {
				// fmt.Printf("receive.1: input send %v\n", n)
				q.buf.Push(n)
			} else {
				inputClosed = true
				myInput = nil
				// fmt.Printf("receive.2 input closed\n")
			}
			common.RandomSleep()

		case op := <-q.request:

			switch op {

			case opQuitDequeue:
				// esce smettendo subito la lettura dell'input (che però non chiude)
				// ma continua a restituire gli elementi già letti dall'input fino all'esaurimento della coda
				inputClosed = true
				myInput = nil
				// fmt.Printf("opQuitDequeue\n")

			case opQuitNext:
				// esce smettendo subito la lettura dell'input (che però non chiude)
				// ma deve attendere la prossima richiesta di lettura prima di chiudere il response channel
				inputClosed = true
				myInput = nil
				closingResp = true
				// fmt.Printf("opQuitNext closing ...\n")

			case opGetItem:
				if closingResp {
					// fmt.Printf("opGetItem: prec opQuitNext close now\n")
					close(q.response)
					return
				}
				L := q.buf.Len()
				// fmt.Printf("opGetItem: LEN = %d - CAP = %d\n", L, q.buf.Cap())
				if L > 0 {
					// FIFO
					n := q.buf.Pop()
					// fmt.Printf("opGetItem: BUF -> %v\n", n)
					q.response <- n
				} else { // L==0
					if inputClosed {
						// fmt.Printf("opGetItem: input closed: closing response channel\n")
						close(q.response)
						return
					}
					// fmt.Printf("opGetItem: waiting for input\n")
					n, ok := <-input
					if ok {
						// fmt.Printf("opGetItem: INPUT -> %v\n", n)
						q.response <- n
					} else {
						// fmt.Printf("opGetItem: INPUT CLOSED: CLOSING RESP\n")
						close(q.response)
						return
					}
				}
			}
		}
	}
}

func (q *InfBufChan) quitNext() {
	q.request <- opQuitNext
}

func (q *InfBufChan) quitDequeue() {
	q.request <- opQuitDequeue
}

func (q *InfBufChan) consume() <-chan interface{} {
	output := make(chan interface{})
	go func() {
		for {
			// ask for an item
			q.request <- opGetItem
			// wait for the response
			item, ok := <-q.response
			if ok {
				output <- item
			} else {
				close(output)
				return
			}
		}
	}()
	return output
}
