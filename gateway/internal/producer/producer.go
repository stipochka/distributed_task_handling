package producer

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Producer interface {
	SendMessage(ctx context.Context, message []byte) error
}

type NatsProducer struct {
	stream jetstream.JetStream
}

func NewNatsProducer(natsAddr string) (*NatsProducer, error) {
	nc, err := nats.Connect(natsAddr)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	_, err = js.CreateStream(context.Background(), jetstream.StreamConfig{
		Name:     "TASKS",
		Subjects: []string{"tasks.*"},
	})
	if err != nil {
		return nil, err
	}

	return &NatsProducer{
		stream: js,
	}, nil
}

func (np *NatsProducer) SendMessage(ctx context.Context, msg []byte) error {
	_, err := np.stream.Publish(ctx, "tasks.create", msg)
	return err
}
