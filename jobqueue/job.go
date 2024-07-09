package jobqueue

type Job interface {
	Run() JobResult
}

type JobResult struct {
	Successful   bool
	ErrorMessage string
}

type JobWrapper struct {
	Job
	tries int
}
