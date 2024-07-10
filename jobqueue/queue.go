package jobqueue

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

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
func (q *Queue) AddJob(j Job) {
	jw := JobWrapper{Job: j, tries: 0, id: uuid.NewString()}
	q.jobs <- jw
	fmt.Printf("job %s created\n", jw.id)

}

func (q *Queue) worker(jobs <-chan JobWrapper) {
	for job := range jobs {
		fmt.Printf("job %s accepted\n", job.id)

		c, cancel := context.WithTimeout(context.Background(), q.Timeout)
		r := make(chan JobResult)
		go func() {
			job.tries++
			r <- job.Run()
		}()
		select {
		case <-c.Done():
			fmt.Printf("job %s timed out\n", job.id)

			if job.tries < q.Attempts {
				fmt.Printf("job %s retried\n", job.id)
				q.jobs <- job
			}
		case result := <-r:
			if result.Successful {
				fmt.Printf("job %s completed successfully\n", job.id)

			} else {
				fmt.Printf("job %s errored out: %s\n", job.id, result.ErrorMessage)

				if job.tries < q.Attempts {
					fmt.Printf("job %s retried\n", job.id)

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
