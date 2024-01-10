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
	ActionUpdate = "update"
)

func NewProfileService(repo db.ProfileRepo, broker events.EventProducer) *ProfileService {
	broker.CreateTopics(
		logger.WithCtx(context.Background(), logger.Get()),
		events.TopicConfig{ReplicationFactor: 1, Partitions: 1},
		GetTopicName(ActionCreate),
		GetTopicName(ActionDelete),
		GetTopicName(ActionUpdate),
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

func (s *ProfileService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	topic := GetTopicName(ActionDelete)
	err = s.broker.Publish(ctx, topic, map[string]string{"id": id}, map[string]string{})
	if err != nil {
		return err
	}
	return nil
}

func (s *ProfileService) Update(ctx context.Context, profileId string, update client.ProfileUpdateRequest) error {
	fmt.Printf("update: %v\n", update)
	err := s.repo.UpdateByIdentityID(ctx, profileId, update)
	if err != nil {
		return err
	}
	topic := GetTopicName(ActionUpdate)
	err = s.broker.Publish(ctx, topic, update, map[string]string{})
	if err != nil {
		return err
	}
	return nil
}

func (s *ProfileService) GetByIdentityID(ctx context.Context, id string) (client.ProfileResponse, error) {
	return s.repo.GetByIdentityID(ctx, id)
}

func GetTopicName(action string) string {
	return fmt.Sprintf("general_%s_%s", TopicProfiles, action)
}
