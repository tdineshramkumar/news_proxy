/*
	pacakge request defines the ExecuteTask, SerialRequest and ParallelRequest functions.
	They are utilities that help in different implementations of server.
*/
package request

import (
	"time"
)

// data is an internal data type used for communication of data between goroutines in ExecuteTask
// through a channel.
type data struct {
	value interface{}
	err   error
}

// TimeoutError implements the error interface and denotes the Timeout of ExecuteTask
type TimeoutError string

func (err TimeoutError) Error() string {
	return string(err)
}

// ExecuteTask executes the Task() function and returns the data the function returns.
// However If execution time exceeds the timeout duration, it executes the Task() function again
// and waits for any result till another timeout and this repeats for a maximum of NumRetries.
// If no result is obtained after NumRetries timeouts, it outputs TimeoutError
func ExecuteTask(Task func() (interface{}, error), Timeout time.Duration, NumRetries int) (interface{}, error) {
	ch := make(chan data, NumRetries)
	for run := 0; run < NumRetries; run++ {
		go func() {
			value, err := Task()
			ch <- data{value: value, err: err}
		}()
		select {
		case data := <-ch:
			return data.value, data.err
		case <-time.After(Timeout):
		}
	}
	return nil, TimeoutError("All Requests Timed Out.")
}
