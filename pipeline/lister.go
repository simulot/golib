package pipeline

import (
	"fmt"
)

type Stringer interface {
	String() string
}

func ListerOperator() Operator {
	return func(in, out chan interface{}) {
		for item := range in {
			if s, ok := item.(Stringer); ok {
				fmt.Println(s.String())
			} else {
				fmt.Printf("%v\n", item)
			}
			out <- item
		}
	}
}
