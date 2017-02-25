package pipeline

import (
	"sync"
)

type Runner interface {
	Run(interface{}, chan error) chan interface{}
}

type Operator func(chan interface{}, chan interface{})

func (o Operator) Run(in chan interface{}) chan interface{} {
	out := make(chan interface{})
	go func() {
		o(in, out)
		close(out)
	}()
	return out
}

type Flow []Operator

func NewFlow(ops ...Operator) Flow {
	return Flow(ops)
}

func (f Flow) Run(in chan interface{}) chan interface{} {
	for _, o := range f {
		if o != nil {
			in = o.Run(in)
		}
	}
	return in
}

type ParallelFlow struct {
	Flow Flow
	n    int
}

func NewParallelFlow(n int, ops ...Operator) Operator {
	w := ParallelFlow{
		Flow: Flow(ops),
		n:    n,
	}
	return w.Run
}

func (w *ParallelFlow) Run(in chan interface{}, out chan interface{}) {
	wg := sync.WaitGroup{}
	wg.Add(w.n)
	for i := 0; i < w.n; i++ {
		go func() {
			localOut := w.Flow.Run(in)
			for item := range localOut {
				out <- item
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
