package service

import (
	"context"
	"pr-assignment/internal/adapter/out/repository"
	"pr-assignment/internal/model"

	"github.com/google/uuid"
)

type UserService struct {
	userRepository *repository.UserRepository
	teamRepository *repository.TeamRepository
}

func NewUserService(r *repository.UserRepository, t *repository.TeamRepository) *UserService {
	return &UserService{r, t}
}

func (s *UserService) SetUserActive(ctx context.Context, userId string, isActive bool) (*model.User, error) {
	user, err := s.userRepository.UpdateUserStatus(ctx, userId, isActive)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) AddTeam(ctx context.Context, team model.Team) error {
	res, err := s.teamRepository.Exists(ctx, team.TeamName)
	if err != nil {
		return err
	}
	if res {
		return model.NewError(model.TEAM_EXISTS, "%s already exists", team.TeamName)
	}

	teamId := uuid.New()
	err = s.teamRepository.AddTeam(ctx, team, teamId)
	if err != nil {
		return err
	}

	err = s.userRepository.AddTeam(ctx, team, teamId)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetTeam(ctx context.Context, teamName string) (*model.Team, error) {
	res, err := s.teamRepository.Exists(ctx, teamName)
	if err != nil {
		return nil, err
	}
	if !res {
		return nil, model.NewError(model.NOT_FOUND, "%s team table not found", teamName)
	}

	team, err := s.userRepository.GetTeam(ctx, teamName)
	if err != nil {
		return nil, model.NewError(model.NOT_FOUND, "%s not found", teamName)
	}

	return team, nil
}
