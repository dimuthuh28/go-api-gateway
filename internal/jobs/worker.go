package jobs

import "fmt"

func ProcessJob(job Job) {
	// Replace with real logic
	fmt.Println("Processing job:", job.Request)
}

func StartWorkers(jobQueue chan Job, workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func() {
			for job := range jobQueue {
				ProcessJob(job)
			}
		}()
	}
}
