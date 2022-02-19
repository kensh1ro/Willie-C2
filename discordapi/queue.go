package discordapi

import "fmt"

type Queue []interface{}

func (Q *Queue) Push(x interface{}) {
	*Q = append(*Q, x)
}

func (Q *Queue) Pop() interface{} {
	h := *Q
	fmt.Printf("POP: %s\n", h)
	var el interface{}
	l := len(h)
	el, *Q = h[0], h[1:l]
	return el
}

func NewQueue() *Queue {
	return &Queue{}
}
