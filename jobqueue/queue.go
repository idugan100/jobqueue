package jobqueue

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type priority int

const HIGH = priority(3)
const MEDIUM = priority(2)
const LOW = priority(1)

type Queue struct {
	Workers int
	Retries int
	Timeout time.Duration
	Size    int
	wg      sync.WaitGroup
	jobs    chan Job
}

func NewQueue() *Queue {
	q := &Queue{Workers: 3, Retries: 1, Timeout: 1 * time.Second, Size: 5}
	q.jobs = make(chan Job, q.Size)

	for range q.Workers {
		q.wg.Add(1)
		fmt.Println("worker created")
		go q.worker(q.jobs)
	}
	return q
}
func (q *Queue) AddJob(j Job, p priority) {
	fmt.Println("job created")
	q.jobs <- j
}

func (q *Queue) worker(jobs <-chan Job) {
	for job := range jobs {
		fmt.Println("job accepted")
		c, cancel := context.WithTimeout(context.Background(), q.Timeout)
		r := make(chan JobResult)
		go func() {
			r <- job.Run()
		}()
		select {
		case <-c.Done():
			fmt.Println("TIMEOUT")
		case result := <-r:
			if result.Successful {
				fmt.Println("job compeleted successfully")
			} else {
				fmt.Printf("error completing job %s", result.ErrorMessage)
			}
		}
		cancel()
	}
	q.wg.Done()

}

func (q *Queue) Stop() {
	close(q.jobs)
	q.wg.Wait()
}
