package pool

import (
	"context"
	"fmt"
	"sync"
	"time"
	"worker/internal/models"
	"worker/internal/storage"

	"github.com/sirupsen/logrus"
)

const (
	chanCap = 10
)

type WorkerPool interface {
	Start(ctx context.Context)
	SubmitTask(task models.Task)
	Close()
}

type Pool struct {
	workersCount int
	taskCh       chan models.Task
	repo         storage.Storage
	wg           *sync.WaitGroup
}

func NewPool(workers int, repo storage.Storage) *Pool {
	return &Pool{
		workersCount: workers,
		taskCh:       make(chan models.Task, chanCap),
		repo:         repo,
		wg:           &sync.WaitGroup{},
	}
}

func (p *Pool) Start(ctx context.Context) {
	for id := 0; id < p.workersCount; id++ {
		id := id
		p.wg.Add(1)
		go p.work(ctx, id)
	}
}

func (p *Pool) SubmitTask(task models.Task) {
	p.taskCh <- task
}

func (p *Pool) work(ctx context.Context, id int) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			logrus.WithField("worker_id", id).Info("context expired")
			return
		case task, ok := <-p.taskCh:
			if !ok {
				logrus.Infof("task chan is closed worker_id=%d", id)
			}
			logrus.
				WithFields(logrus.Fields{
					"worker_id": id,
					"task_id":   task.TaskID,
				}).
				Info("worker received task")

			var err error

			var status, description string

			switch task.Type {
			case "resize image":
				err = p.ResizeImage(&task)
				status = "processed"
				description = fmt.Sprintf("task %s successfully completed", task.Type)
			default:
				logrus.Infof("received unknown task type, type: %s, worker_id=%d", task.Type, id)
				status = "unknown type"
				description = fmt.Sprintf("task %s is unknown", task.Type)
			}

			if err != nil {
				logrus.WithError(err).Errorf("failed to proceed task, worker_id=%d", id)
				status = "failed with error"
				description = fmt.Sprintf("failed with error %s", err.Error())
			}

			result := models.Result{
				TaskID:      task.TaskID,
				Type:        task.Type,
				Status:      status,
				Description: description,
			}

			if err := p.repo.SaveResult(ctx, result); err != nil {
				logrus.WithError(err).Errorf("failed to save task result, worker_id=%d", id)
				continue
			}

		}
	}
}

func (p *Pool) Close() {
	close(p.taskCh)
	p.wg.Wait()
}

func (p *Pool) ResizeImage(task *models.Task) error {
	// logic of resize image
	time.Sleep(time.Millisecond * 500) // imitating work
	logrus.Infof("task %s resized successfully", task.TaskID)
	return nil
}
