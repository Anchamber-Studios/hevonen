package services

import (
	"context"
	"fmt"

	"github.com/anchamber-studios/hevonen/lib/events"
	"github.com/anchamber-studios/hevonen/lib/logger"
	"github.com/anchamber-studios/hevonen/services/general/profile/client"
	"github.com/anchamber-studios/hevonen/services/general/profile/db"
)

type ProfileService struct {
	repo   db.ProfileRepo
	broker events.EventProducer
}

const (
	TopicProfiles = "profiles"
)

const (
	Version = "1.0.0"
)

const (
	ActionCreate = "create"
	ActionDelete = "delete"
)

func NewProfileService(repo db.ProfileRepo, broker events.EventProducer) *ProfileService {
	broker.CreateTopics(
		logger.WithCtx(context.Background(), logger.Get()),
		events.TopicConfig{ReplicationFactor: 1, Partitions: 1},
		GetTopicName(ActionCreate),
		GetTopicName(ActionDelete),
	)
	return &ProfileService{
		repo:   repo,
		broker: broker,
	}
}

func (s *ProfileService) Create(ctx context.Context, profile client.ProfileCreateRequest) (string, error) {
	id, err := s.repo.Create(ctx, profile)
	if err != nil {
		return "", err
	}
	topic := GetTopicName(ActionCreate)
	err = s.broker.Publish(ctx, topic, profile, map[string]string{})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *ProfileService) Get(ctx context.Context, id string) (client.ProfileResponse, error) {
	return s.repo.Get(ctx, id)
}

func GetTopicName(action string) string {
	return fmt.Sprintf("general_%s_%s", TopicProfiles, action)
}
