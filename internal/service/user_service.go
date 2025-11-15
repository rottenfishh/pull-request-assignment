package service

import (
	"context"
	"fmt"
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
	fmt.Println("service ", teamName)
	teamId, err := s.teamRepository.GetTeamId(ctx, teamName)

	if err != nil {
		fmt.Println("error getting team id teamRepo")
		return nil, err
	}

	team, err := s.userRepository.GetTeam(ctx, teamId)
	if err != nil {
		fmt.Println("error getting team userRepo")
		fmt.Println(err)
		return nil, model.NewError(model.NOT_FOUND, "%s not found", teamName)
	}
	team.TeamName = teamName
	return team, nil
}

func (s *PullRequestService) GetActiveTeammatesByUser(ctx context.Context, userId string) ([]string, error) {
	teamId, err := s.userRepository.GetTeamNameByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	userIds, err := s.userRepository.GetActiveUsersByTeam(ctx, teamId)
	if err != nil {
		return nil, err
	}

	return userIds, nil
}
