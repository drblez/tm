package tm

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Tm struct {
	Batch     int
	Defer     int
	WaitGroup sync.WaitGroup
	JobQueue  chan Job
	Workers   map[string]WorkerIteraface
}

func (tm *Tm) Add() {
	tm.WaitGroup.Add(1)
}

func (tm *Tm) Wait() {
	close(tm.JobQueue)
	tm.WaitGroup.Wait()
}

func (tm *Tm) Done() {
	tm.WaitGroup.Done()
}

func Create(q int) Tm {
	tm := Tm{
		Batch:     30,
		Defer:     500,
		WaitGroup: sync.WaitGroup{},
		JobQueue:  make(chan Job, q),
		Workers:   make(map[string]WorkerIteraface),
	}
	return tm
}

func (tm *Tm) Register(w WorkerIteraface) {
	w.Init()
	tm.Add()
	go func() {
		defer tm.Done()
		for {
			j := make(chan Job)

			tm.Add()
			go func() {
				defer tm.Done()
				a := make([]interface{}, 0)
				for jj := range j {
					a = append(a, jj.Args[0])
				}

				if len(a) > 0 {
					job := Job{w.GetType(), a}
					w.Do(&job)
				}

			}()

		loop:
			for r := 0; r < tm.Batch; {
				select {

				case <-time.After(time.Duration(tm.Defer) * time.Millisecond):
					if r > 0 {
						break loop
					}

				case job, open := <-w.GetJobs():
					if !open {
						close(j)
						return
					}
					j <- job
					r++
				}
			}

			close(j)
		}

	}()

	tm.Workers[w.GetType()] = w
}

func (tm *Tm) Dispatch() {
	tm.Add()
	go func() {
		defer tm.Done()
		for job := range tm.JobQueue {
			if _, ok := tm.Workers[job.Type]; !ok {
				log.Println(fmt.Sprintf("Invalid job type %s", job.Type))
				continue
			}
			tm.Workers[job.Type].GetJobs() <- job
		}

		for k, _ := range tm.Workers {
			close(tm.Workers[k].GetJobs())
		}

	}()
}

type WorkerIteraface interface {
	Init()

	GetType() string
	GetJobs() chan Job
	Do(job *Job)
}

type Worker struct {
	Type string
	Jobs chan Job
}

func (w *Worker) GetJobs() chan Job {
	return w.Jobs
}

func (w *Worker) Do(job *Job) {
	fmt.Println("Custom worker 'Do' method not implemented")
}

func (w *Worker) Init() {
	w.Jobs = make(chan Job)
}

func (w *Worker) GetType() string {
	return w.Type
}

type Job struct {
	Type string
	Args []interface{}
}
