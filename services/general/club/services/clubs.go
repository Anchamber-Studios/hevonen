package services

import (
	"context"

	"github.com/anchamber-studios/hevonen/services/club/client"
	"github.com/anchamber-studios/hevonen/services/club/db"
)

type ClubService struct {
	repo db.ClubRepo
}

func NewClubService(repo db.ClubRepo) *ClubService {
	return &ClubService{repo: repo}
}

func (s *ClubService) List(ctx context.Context) ([]client.Club, error) {
	return s.repo.List(ctx)
}

func (s *ClubService) ListForIdentity(ctx context.Context, identity string) ([]client.ClubMember, error) {
	return s.repo.ListForIdentity(ctx, identity)
}

func (s *ClubService) Create(ctx context.Context, club client.ClubCreate) (string, error) {
	return s.repo.Create(ctx, club)
}

type MemberService struct {
	repo db.MemberRepo
}

func NewMemberService(repo db.MemberRepo) *MemberService {
	return &MemberService{repo: repo}
}

func (s *MemberService) List(ctx context.Context) ([]client.Member, error) {
	return s.repo.List(ctx)
}

func (s *MemberService) ListForClub(ctx context.Context, cId string) ([]client.Member, error) {
	return s.repo.ListForClub(ctx, cId)
}

func (s *MemberService) Create(ctx context.Context, member client.MemberCreate) (string, error) {
	return s.repo.Create(ctx, member)
}
