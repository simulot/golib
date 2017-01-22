package pipeline

// CounterOperator is counter operator. This operator count the number of item going through its inpout
func CounterOperator() Operator {
	return func(in chan interface{}, out chan interface{}) {
		c := 0
		for _ = range in {
			c++
		}
		out <- c
	}
}
