package util

// Queue ...
type Queue interface {
	Push(n interface{})
	Pop() interface{}
	Len() int
	Cap() int
}

// NewCircularFifoQueue returns a new baseCircularFifoQueue
// If needed, the new capacity will be: L_new = 2*L_old + 1
func NewCircularFifoQueue(initCap int) Queue {
	fn := func(L int) int { return 2*L + 1 }
	return NewBaseCircularFifoQueue(initCap, fn)
}

// NewLinearExtendCircularFifoQueue returns a new baseCircularFifoQueue
// with the given initial size.
// If needed, the new capacity will be: L_new = L_old + size
func NewLinearExtendCircularFifoQueue(size int) Queue {
	fn := func(L int) int { return L + size }
	return NewBaseCircularFifoQueue(size, fn)
}

// NewBaseCircularFifoQueue returns a new baseCircularFifoQueue
// with the given ExtendCapFunc function.
func NewBaseCircularFifoQueue(initCap int, fnExtendCap ExtendCapFunc) Queue {
	if initCap <= 0 {
		panic("initCap must be > 0")
	}
	return &baseCircularFifoQueue{
		nodes:       make([]interface{}, initCap),
		fnExtendCap: fnExtendCap,
	}
}

// ExtendCapFunc ...
type ExtendCapFunc func(int) int

// baseCircularFifoQueue is a basic FIFO baseCircularFifoQueue based on a circular list that resizes as needed.
type baseCircularFifoQueue struct {
	nodes       []interface{}
	fnExtendCap ExtendCapFunc
	head        int
	tail        int
	count       int
}

// Push adds a node to the baseCircularFifoQueue.
func (q *baseCircularFifoQueue) Push(n interface{}) {
	if q.head == q.tail && q.count > 0 {
		Lold := len(q.nodes)        // old capacity
		Lnew := q.fnExtendCap(Lold) // new capacity
		/*
			if Lnew < Lold {
				panic(fmt.Sprintf("fnExtendCap(%d) -> %d : new cap must be greater the old cap!", Lold, Lnew))
			}
		*/
		nodes := make([]interface{}, Lnew)
		copy(nodes, q.nodes[q.head:])
		copy(nodes[Lold-q.head:], q.nodes[:q.head])
		q.head = 0
		q.tail = Lold
		q.nodes = nodes
	}
	q.nodes[q.tail] = n
	q.tail = (q.tail + 1) % len(q.nodes)
	q.count++
}

// Pop removes and returns a node from the baseCircularFifoQueue in first to last order.
func (q *baseCircularFifoQueue) Pop() interface{} {
	if q.count == 0 {
		return nil
	}
	node := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.count--
	return node
}

func (q *baseCircularFifoQueue) Len() int {
	return q.count
}

func (q *baseCircularFifoQueue) Cap() int {
	return len(q.nodes)
}

/*
func (q *baseCircularFifoQueue) String() string {
	return fmt.Sprintf("baseCircularFifoQueue{len:%d, cap:%d}", q.Len(), q.Cap())
}
*/
