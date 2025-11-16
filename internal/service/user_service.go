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

func (s *UserService) SetUserActive(ctx context.Context, userID string, isActive bool) (*model.User, error) {
	user, err := s.userRepository.UpdateUserStatus(ctx, userID, isActive)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) AddTeam(ctx context.Context, team model.Team) error {
	res, _ := s.teamRepository.Exists(ctx, team.TeamName)

	if res {
		return model.NewError(model.TeamExists, "%s already exists", team.TeamName)
	}

	teamID := uuid.New()
	err := s.teamRepository.AddTeam(ctx, team, teamID)
	if err != nil {
		return err
	}

	err = s.userRepository.AddTeam(ctx, team, teamID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetTeam(ctx context.Context, teamName string) (*model.Team, error) {
	teamID, err := s.teamRepository.GetTeamID(ctx, teamName)

	if err != nil {
		fmt.Println("error getting team id teamRepo")
		return nil, err
	}

	team, err := s.userRepository.GetTeam(ctx, teamID)
	if err != nil {
		fmt.Println("error getting team userRepo")
		fmt.Println(err)
		return nil, model.NewError(model.NotFound, "%s not found", teamName)
	}
	team.TeamName = teamName
	return team, nil
}

func (s *UserService) GetActiveTeammatesByUser(ctx context.Context, userID string) ([]string, error) {
	teamID, err := s.userRepository.GetTeamNameByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userIDs, err := s.userRepository.GetActiveUsersByTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}

	return userIDs, nil
}

func (s *UserService) KillTeam(ctx context.Context, teamName string) (*model.Team, error) {
	teamID, err := s.teamRepository.GetTeamID(ctx, teamName)
	if err != nil {
		return nil, err
	}

	teamMembers, err := s.userRepository.GetActiveUsersByTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}

	for _, member := range teamMembers {
		_, err = s.userRepository.UpdateUserStatus(ctx, member, false)
		if err != nil {
			return nil, err
		}
	}

	team, err := s.userRepository.GetTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}

	return team, nil
}
