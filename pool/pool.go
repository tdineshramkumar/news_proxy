package pool

import (
	"fmt"
)

const (
	BUF_FACTOR = 2
)

type request struct {
	task     func() (interface{}, error)
	quit     <-chan bool
	response chan<- interface{}
}

type Pool struct {
	requests chan request
}

func (p *Pool) function() {
	for request := range p.requests {
		select {
		case <-request.quit:
			close(request.response)
		default:
			if result, err := request.task(); err == nil {
				// Make sure response channel is buffered.
				request.response <- result
			} else {
				fmt.Println("ERROR: Task returned error in pool function [", err, "]")
				close(request.response)
			}
		}
	}
}
func New(size int) *Pool {
	pool := &Pool{
		requests: make(chan request, size*BUF_FACTOR),
	}
	for i := 0; i < size; i++ {
		go pool.function()
	}
	return pool
}

func (p *Pool) Execute(task func() (interface{}, error), quit <-chan bool) interface{} {
	response := make(chan interface{}, 1)
	p.requests <- request{task: task, quit: quit, response: response}
	return <-response
}
