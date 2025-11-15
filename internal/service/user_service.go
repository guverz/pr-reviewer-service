package service

import (
	"context"

	"github.com/guverz/pr-reviewer-service/internal/domain"
	"github.com/guverz/pr-reviewer-service/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewUserService(userRepo repository.UserRepository, teamRepo repository.TeamRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

// SetActive устанавливает флаг активности пользователя и синхронизирует данные в команде
func (s *UserService) SetActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	// Получаем пользователя для получения teamName
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.NewDomainError(domain.ErrorCodeNotFound, "user not found")
	}

	// Обновляем статус активности в UserRepository
	updatedUser, err := s.userRepo.SetActive(ctx, userID, isActive)
	if err != nil {
		return nil, domain.NewDomainError(domain.ErrorCodeNotFound, "user not found")
	}

	// Синхронизируем данные в TeamRepository
	if err := s.teamRepo.UpdateMember(ctx, user.TeamName, userID, isActive); err != nil {
		// Возвращаем ошибку для отладки - если команда или участник не найдены,
		// это означает проблему синхронизации данных
		return nil, domain.NewDomainError(domain.ErrorCodeNotFound, "failed to update team member: "+err.Error())
	}

	return updatedUser, nil
}

// GetUser получает пользователя по ID
func (s *UserService) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.NewDomainError(domain.ErrorCodeNotFound, "user not found")
	}
	return user, nil
}
