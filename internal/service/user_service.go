package service

import (
	"pr-assignment/internal/adapter/out/repository"
	"pr-assignment/internal/model"
)

type UserService struct {
	userRepository *repository.UserRepository
	teamRepository *repository.TeamRepository
}

func NewUserService(r *repository.UserRepository, t *repository.TeamRepository) *UserService {
	r.Init()
	t.Init()
	return &UserService{r, t}
}

func (s *UserService) SetUserActive(userId string, isActive bool) error {
	err := s.userRepository.UpdateUserStatus(userId, isActive)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) AddTeam(team model.Team) error {
	res, err := s.teamRepository.Exists(team.TeamName)
	if err != nil {
		return err
	}
	if res {
		return model.NewError(model.TEAM_EXISTS, "%s already exists", team.TeamName)
	}
	err = s.userRepository.AddTeam(team)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetTeam(teamName string) (*model.Team, error) {
	res, err := s.teamRepository.Exists(teamName)
	if err != nil {
		return nil, err
	}
	if !res {
		return nil, model.NewError(model.NOT_FOUND, "%s not found", teamName)
	}

	team, err := s.userRepository.GetTeam(teamName)
	if err != nil {
		return nil, model.NewError(model.NOT_FOUND, "%s not found", teamName)
	}

	return team, nil
}
