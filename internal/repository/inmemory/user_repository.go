package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/guverz/pr-reviewer-service/internal/domain"
)

type UserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *UserRepository) UpsertTeamMembers(ctx context.Context, teamName string, members []domain.TeamMember) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, member := range members {
		user := &domain.User{
			ID:       member.UserID,
			Username: member.Username,
			TeamName: teamName,
			IsActive: member.IsActive,
		}
		r.users[member.UserID] = user
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}

	// Возвращаем копию
	userCopy := *user
	return &userCopy, nil
}

func (r *UserRepository) SetActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}

	// Обновляем флаг активности
	user.IsActive = isActive

	// Возвращаем копию
	userCopy := *user
	return &userCopy, nil
}

func (r *UserRepository) ListByTeam(ctx context.Context, teamName string, onlyActive bool) ([]domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]domain.User, 0)
	for _, user := range r.users {
		if user.TeamName == teamName {
			if onlyActive && !user.IsActive {
				continue
			}
			users = append(users, *user)
		}
	}

	return users, nil
}



