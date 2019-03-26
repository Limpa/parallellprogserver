package sandbox

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os/exec"
	"sync"
	"time"
)

var (
	q jobqueue
)

type job struct {
	File        *multipart.File
	TimeCreated time.Time
	Status      string
	Name        string
}

type jobqueue struct {
	Queue         []job
	Channel       chan *multipart.File
	Lock          sync.Mutex
	OnQueueChange func()
}

func init() {
	q = jobqueue{make([]job, 0), make(chan *multipart.File, 100), sync.Mutex{}, nil}
}

func compileAndExecute(f *multipart.File) error {
	fb, err := ioutil.ReadAll(*f)
	if err != nil {
		return err
	}
	ioutil.WriteFile("/tmp/tmp.go", fb, 0644)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if exec.CommandContext(ctx, "go", "run", "/tmp/tmp.go").Run(); err != nil {
		return err
	}

	return nil
}

// Run starts the sandbox environment
func Run(onQueueChange func()) {
	q.OnQueueChange = onQueueChange
	for j := range q.Channel {
		q.Lock.Lock()
		job := q.Queue[0]
		q.Queue[0].Status = "Executing"
		onQueueChange()
		q.Queue = q.Queue[1:]
		q.Lock.Unlock()
		fmt.Printf("Executing file created at %s\n", job.TimeCreated.String())
		err := compileAndExecute(j)
		if err != nil {
			fmt.Println(err)
		}
		job.Status = "Executed"
		onQueueChange()
	}
}

// GetQueue returns a json representation of the queue as a byte array
func GetQueue() []byte {
	var buffer bytes.Buffer
	buffer.WriteString(`{"queue":[`)
	for i, job := range q.Queue {
		buffer.WriteString(fmt.Sprintf(`{"name":"%s","status":"%s","qtime":"%s"}`, job.Name, job.Status, time.Since(job.TimeCreated).String()))
		if i < len(q.Queue)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString(`]}`)
	return buffer.Bytes()
}

// CreateJob creates a job for the sandbox
func CreateJob(f *multipart.File, name string) {
	q.Lock.Lock()
	q.Queue = append(q.Queue, job{f, time.Now(), "Waiting", name})
	q.OnQueueChange()
	q.Lock.Unlock()
	q.Channel <- f
}
