package consumer

import (
	"context"
	"encoding/json"
	"worker/internal/models"
	"worker/internal/pool"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	js   nats.JetStream
	pool pool.WorkerPool
}

func NewConsumer(addr string, pool pool.WorkerPool) (*Consumer, error) {
	nc, err := nats.Connect(addr)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return &Consumer{
		js:   js,
		pool: pool,
	}, nil
}

func (c *Consumer) Run(ctx context.Context, streamName string) error {
	sub, err := c.js.Subscribe(streamName, func(msg *nats.Msg) {
		task := models.Task{}

		if err := json.Unmarshal(msg.Data, &task); err != nil {
			logrus.WithError(err).Error("failed to parse task")
			_ = msg.Ack()
			return
		}

		c.pool.SubmitTask(task)
		_ = msg.Ack()
	}, nats.Durable("worker1"), nats.ManualAck())

	if err != nil {
		return err
	}

	defer sub.Unsubscribe()
	<-ctx.Done()

	logrus.Info("stopping consumer...")
	if err := sub.Unsubscribe(); err != nil {
		logrus.WithError(err).Error("failed to unsubscribe")
	}

	c.pool.Close()
	return nil
}
