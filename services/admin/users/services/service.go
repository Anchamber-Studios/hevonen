package services

import (
	"context"
	"fmt"

	"github.com/anchamber-studios/hevonen/lib/events"
	"github.com/anchamber-studios/hevonen/lib/logger"
	"github.com/anchamber-studios/hevonen/services/admin/users/client"
	"github.com/anchamber-studios/hevonen/services/admin/users/db"
)

type UserService struct {
	repo   db.UserRepo
	broker events.EventProducer
}

const (
	Version = "1.0.0"
)

const (
	ActionCreate         = "create"
	ActionUpdateEmail    = "update_email"
	ActionChangePassword = "change_password"
	ActionDelete         = "delete"
	ActionConfirmEmail   = "confirm_email"
	ActionLogin          = "login"
	ActionLogout         = "logout"
	ActionDeactivate     = "deactivate"
	ActionReactivate     = "reactivate"
)

func NewUserService(repo db.UserRepo, broker events.EventProducer) *UserService {
	return &UserService{
		repo:   repo,
		broker: broker,
	}
}

func (s *UserService) Create(ctx context.Context, user client.UserCreate) (string, error) {
	id, err := s.repo.Create(ctx, user)
	if err != nil {
		return "", err
	}
	topic := getTopicName(ActionCreate, id)
	err = s.broker.Publish(ctx, topic, user)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *UserService) Login(ctx context.Context, login client.UserLogin) (client.User, error) {
	user, err := s.repo.Login(ctx, login)
	if err != nil {
		return client.User{}, err
	}
	topic := getTopicName(ActionLogin, user.ID)
	msgCtx := logger.WithCtx(context.Background(), logger.FromContext(ctx))
	err = s.broker.Publish(msgCtx, topic, user)
	if err != nil {
		return client.User{}, err
	}
	return user, nil
}

func (s *UserService) List(ctx context.Context) ([]client.User, error) {
	return s.repo.List(ctx)
}

func (s *UserService) Get(ctx context.Context, id string) (client.User, error) {
	return s.repo.Get(ctx, id)
}

func getTopicName(action string, id string) string {
	return fmt.Sprintf("users_%s_%s_%s", Version, action, id)
}
