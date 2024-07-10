package main

import (
	"fmt"
	"time"

	"github.com/idugan100/jobqueue/jobqueue"
)

type myjob struct{}

func (m myjob) Run() jobqueue.JobResult {
	time.Sleep(5 * time.Second)
	_, err := fmt.Println("hello world!")
	if err != nil {
		return jobqueue.JobResult{Successful: false, ErrorMessage: err.Error()}
	}
	return jobqueue.JobResult{Successful: true}
}

type quickjob struct{}

func (q quickjob) Run() jobqueue.JobResult {
	_, err := fmt.Println("hello world!")
	if err != nil {
		return jobqueue.JobResult{Successful: false, ErrorMessage: err.Error()}
	}
	return jobqueue.JobResult{Successful: true}
}

func main() {
	m := myjob{}
	qj := quickjob{}
	q := jobqueue.NewQueue()

	for range 1 {
		q.AddJob(m)
	}
	for range 2 {
		q.AddJob(qj)
	}
	fmt.Println("hi1")
	fmt.Println("hi2")
	fmt.Println("hi3")
	time.Sleep(time.Second * 30)
	q.Stop()
}
