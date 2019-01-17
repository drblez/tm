<h1>Task manager with mickobatching</h1>

Inspired by https://github.com/kimiby/gotm

<strong>Usage:</strong>
<pre>
package main

import (
	"fmt"
	"github.com/maslennikov-yv/tm"
	"time"
)

type CustomWorker struct {
	tm.Worker
}

func (w CustomWorker) Do(job *tm.Job) {

	for _, a := range job.Args {
		if job, ok := a.(string); ok {
			fmt.Println(job)
		}
	}

	fmt.Println("---")

}

func main() {

	custom := CustomWorker{tm.Worker{Type: "custom"}}

	task_manager := tm.Create(10)
	task_manager.Register(&custom)
	task_manager.Dispatch()

	job := make([]interface{}, 1);
	job[0] = "Custom job"
	task_manager.JobQueue <- tm.Job{"custom", job}
	task_manager.JobQueue <- tm.Job{"custom", job}
	time.Sleep(501* time.Millisecond)
	task_manager.JobQueue <- tm.Job{"custom", job}

	task_manager.Wait()
}

</pre>
