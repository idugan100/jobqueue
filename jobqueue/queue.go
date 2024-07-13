package jobqueue

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Queue struct {
	Workers int
	Timeout time.Duration
	Size    int
	Output  *log.Logger
	wg      sync.WaitGroup
	jobs    chan JobWrapper
}

func NewQueue() *Queue {
	q := &Queue{Workers: 3, Timeout: 1 * time.Second, Size: 5, Output: log.Default()}
	q.jobs = make(chan JobWrapper, q.Size)

	for range q.Workers {
		q.wg.Add(1)
		q.Output.Println("worker created")
		go q.worker()
	}
	return q
}

func (q *Queue) AddJob(j Job) {
	jw := JobWrapper{Job: j, tries: 0, id: uuid.NewString()}
	q.jobs <- jw
	q.Output.Printf("job %s created\n", jw.id)

}

func (q *Queue) worker() {
	var jobs <-chan JobWrapper = q.jobs

	for job := range jobs {
		q.Output.Printf("job %s accepted\n", job.id)

		c, cancel := context.WithTimeout(context.Background(), q.Timeout)
		r := make(chan JobResult)
		go func() {
			job.tries++
			r <- job.Run()
		}()
		select {
		case <-c.Done():
			q.Output.Printf("job %s timed out\n", job.id)
		case result := <-r:
			if result.Successful {
				q.Output.Printf("job %s completed successfully\n", job.id)

			} else {
				q.Output.Printf("job %s errored out: %s\n", job.id, result.ErrorMessage)
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
