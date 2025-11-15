package service

import (
	"context"

	"github.com/guverz/pr-reviewer-service/internal/domain"
	"github.com/guverz/pr-reviewer-service/internal/repository"
)

type TeamService struct {
	teamRepo repository.TeamRepository
	userRepo repository.UserRepository
	txMgr    repository.TransactionManager
}

func NewTeamService(
	teamRepo repository.TeamRepository,
	userRepo repository.UserRepository,
	txMgr repository.TransactionManager,
) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
		txMgr:    txMgr,
	}
}

// CreateTeam создаёт команду и обновляет/создаёт пользователей
func (s *TeamService) CreateTeam(ctx context.Context, team domain.Team) (*domain.Team, error) {
	// Проверяем, существует ли команда
	existing, err := s.teamRepo.GetByName(ctx, team.Name)
	if err == nil && existing != nil {
		return nil, domain.NewDomainError(domain.ErrorCodeTeamExists, "team_name already exists")
	}

	// Создаём команду и пользователей в транзакции
	var createdTeam *domain.Team
	err = s.txMgr.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Создаём команду
		if err := s.teamRepo.Create(txCtx, team); err != nil {
			return err
		}

		// Обновляем/создаём пользователей
		if err := s.userRepo.UpsertTeamMembers(txCtx, team.Name, team.Members); err != nil {
			return err
		}

		// Получаем созданную команду с пользователями
		createdTeam, err = s.teamRepo.GetByName(txCtx, team.Name)
		return err
	})

	if err != nil {
		return nil, err
	}

	return createdTeam, nil
}

// GetTeam получает команду по имени
func (s *TeamService) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	team, err := s.teamRepo.GetByName(ctx, teamName)
	if err != nil {
		return nil, domain.NewDomainError(domain.ErrorCodeNotFound, "team not found")
	}
	return team, nil
}

