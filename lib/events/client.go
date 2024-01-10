package events

import (
	"context"
	"encoding/json"

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

type EventProducer interface {
	TopicExists(ctx context.Context, topic string) (bool, error)
	CreateTopics(ctx context.Context, config TopicConfig, topics ...string) error
	DeleteTopics(ctx context.Context, topics ...string) error
	Publish(ctx context.Context, topic string, msg any, headers map[string]string) error
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

func (c *EventProducerRedpanda) CreateTopics(ctx context.Context, config TopicConfig, topics ...string) error {
	log := logger.FromContext(ctx)
	log.Sugar().Infof("Creating topics '%s'\n", topics)
	resp, err := c.admin.CreateTopics(ctx, int32(config.Partitions), int16(config.ReplicationFactor), config.Configs, topics...)

	if err != nil {
		// log.Sugar().Errorf("Failed to create topic '%s': %v\n", topics, err)
		return err
	}
	for _, r := range resp {
		if r.Err != nil {
			// log.Sugar().Errorf("Failed to create topic '%s': %v\n", topics, r.Err)
			return r.Err
		}
	}
	return nil
}

func (c *EventProducerRedpanda) DeleteTopics(ctx context.Context, topics ...string) error {
	resp, err := c.admin.DeleteTopics(ctx, topics...)
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

func (c *EventProducerRedpanda) Publish(ctx context.Context, topic string, msg any, headers map[string]string) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	log := logger.FromContext(ctx)
	log.Sugar().Infof("Publishing to topic '%s'\n", topic)
	c.client.Produce(ctx, &kgo.Record{
		Topic:   topic,
		Value:   data,
		Headers: MapToKgoHeaders(headers),
	}, func(r *kgo.Record, err error) {
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
	Subscribe(ctx context.Context, handler func(context.Context, *Event) error) error
}

type EventConsumerRedpanda struct {
	client *kgo.Client
}

func NewEventConsumerRedpanda(brokers []string, topics ...string) (*EventConsumerRedpanda, error) {
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
	}, nil
}

func (c *EventConsumerRedpanda) Subscribe(ctx context.Context, handler func(context.Context, []byte, map[string]string) error) error {
	log := logger.FromContext(ctx)
	go func() {
		for {
			fetches := c.client.PollFetches(ctx)
			iter := fetches.RecordIter()
			for !iter.Done() {
				rec := iter.Next()
				err := handler(ctx, rec.Value, FromKgoHeaders(rec.Headers))
				if err != nil {
					log.Sugar().Errorf("error handling event: %v", err)
				}
			}
		}
	}()
	return nil
}

func MapToKgoHeaders(headers map[string]string) []kgo.RecordHeader {
	var recordHeaders []kgo.RecordHeader
	for k, v := range headers {
		recordHeaders = append(recordHeaders, kgo.RecordHeader{
			Key:   k,
			Value: []byte(v),
		})
	}
	return recordHeaders
}

func FromKgoHeaders(headers []kgo.RecordHeader) map[string]string {
	recordHeaders := make(map[string]string) // Initialize the map
	for _, header := range headers {
		recordHeaders[header.Key] = string(header.Value)
	}
	return recordHeaders
}
