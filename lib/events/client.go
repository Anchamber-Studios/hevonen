package events

import (
	"context"
	"encoding/json"
	"log"

	"github.com/anchamber-studios/hevonen/lib/logger"
	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Event interface{}
type TopicConfig struct {
	ReplicationFactor int
	Partitions        int
	Configs           map[string]*string
}

type EventAdmin interface {
	TopicExists(ctx context.Context, topic string) (bool, error)
	CreateTopic(ctx context.Context, topic string, config TopicConfig) error
	DeleteTopic(ctx context.Context, topic string) error
}

type EventProducer interface {
	Publish(ctx context.Context, topic string, msg any) error
}

type EventProducerRedpanda struct {
	admin  *kadm.Client
	client *kgo.Client
}

func NewEventProducerRedpanda(brokers ...string) (*EventProducerRedpanda, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
	)
	if err != nil {
		return nil, err
	}
	admin := kadm.NewClient(client)
	return &EventProducerRedpanda{
		admin:  admin,
		client: client,
	}, nil
}

func (c *EventProducerRedpanda) TopicExists(ctx context.Context, topic string) (bool, error) {
	topicsMetadata, err := c.admin.ListTopics(ctx)
	if err != nil {
		return false, err

	}
	for _, metadata := range topicsMetadata {
		if metadata.Topic == topic {
			return true, nil
		}
	}
	return false, nil
}

func (c *EventProducerRedpanda) CreateTopic(ctx context.Context, topic string, config TopicConfig) error {
	log := logger.FromContext(ctx)
	log.Sugar().Infof("Creating topic '%s'\n", topic)
	resp, err := c.admin.CreateTopic(ctx, int32(config.Partitions), int16(config.ReplicationFactor), config.Configs, topic)
	if err != nil {
		log.Sugar().Errorf("Failed to create topic '%s': %v\n", topic, err)
		return err
	}
	if resp.Err != nil {
		log.Sugar().Errorf("Failed to create topic '%s': %v\n", topic, resp.Err)
		return resp.Err
	}
	return nil
}

func (c *EventProducerRedpanda) DeleteTopic(ctx context.Context, topic string) error {
	resp, err := c.admin.DeleteTopics(ctx, topic)
	if err != nil {
		return err
	}
	for _, r := range resp {
		if r.Err != nil {
			return r.Err
		}
	}
	return nil
}

func (c *EventProducerRedpanda) Publish(ctx context.Context, topic string, msg any) error {
	if exists, err := c.TopicExists(ctx, topic); err != nil && !exists {
		c.CreateTopic(ctx, topic, TopicConfig{ReplicationFactor: 1, Partitions: 1})
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	log := logger.FromContext(ctx)
	log.Sugar().Infof("Publishing to topic '%s'\n", topic)
	c.client.Produce(ctx, &kgo.Record{Topic: topic, Value: data}, func(r *kgo.Record, err error) {
		if err != nil {
			log.Sugar().Errorf("Failed to publish to topic '%s': %v\n", topic, err)
		}
	})
	return nil
}

func (c *EventProducerRedpanda) Close() {
	c.admin.Close()
	c.client.Close()
}

type EventConsumer interface {
	Subscribe(ctx context.Context, topic string, handler func(context.Context, *Event) error) error
}

type EventConsumerRedpanda struct {
	logger *log.Logger
	client *kgo.Client
}

func NewEventConsumerRedpanda(logger *log.Logger, brokers []string, topics ...string) (*EventConsumerRedpanda, error) {
	groupID := uuid.New().String()
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(topics...),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		return nil, err
	}
	return &EventConsumerRedpanda{
		client: client,
		logger: logger,
	}, nil
}

func (c *EventConsumerRedpanda) Subscribe(ctx context.Context, topic string, handler func(context.Context, Event) error) error {
	go func() {
		for {
			fetches := c.client.PollFetches(ctx)
			iter := fetches.RecordIter()
			for !iter.Done() {
				rec := iter.Next()
				var event Event
				err := json.Unmarshal(rec.Value, &event)
				if err != nil {
					c.logger.Printf("error unmarshalling event: %v", err)
				}
				err = handler(ctx, event)
				if err != nil {
					c.logger.Printf("error handling event: %v", err)
				}
			}
		}
	}()
	return nil
}
