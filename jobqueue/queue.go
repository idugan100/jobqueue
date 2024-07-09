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
	Workers  int
	Attempts int
	Timeout  time.Duration
	Size     int
	wg       sync.WaitGroup
	jobs     chan JobWrapper
}

func NewQueue() *Queue {
	q := &Queue{Workers: 3, Attempts: 2, Timeout: 1 * time.Second, Size: 5}
	q.jobs = make(chan JobWrapper, q.Size)

	for range q.Workers {
		q.wg.Add(1)
		fmt.Println("worker created")
		go q.worker(q.jobs)
	}
	return q
}
func (q *Queue) AddJob(j Job, p priority) {
	fmt.Println("job created")
	q.jobs <- JobWrapper{Job: j, tries: 0}
}

func (q *Queue) worker(jobs <-chan JobWrapper) {
	for job := range jobs {
		fmt.Println("job accepted")
		c, cancel := context.WithTimeout(context.Background(), q.Timeout)
		r := make(chan JobResult)
		go func() {
			job.tries++
			r <- job.Run()
		}()
		select {
		case <-c.Done():
			fmt.Println("TIMEOUT")

			if job.tries < q.Attempts {
				fmt.Println("job retried")
				q.jobs <- job
			}
		case result := <-r:
			if result.Successful {
				fmt.Println("job compeleted successfully")
			} else {
				fmt.Printf("error completing job %s\n", result.ErrorMessage)
				if job.tries < q.Attempts {
					fmt.Println("job retried")
					q.jobs <- job
				}
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
