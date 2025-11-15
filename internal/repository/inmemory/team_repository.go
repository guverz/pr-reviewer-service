package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/guverz/pr-reviewer-service/internal/domain"
)

type TeamRepository struct {
	mu    sync.RWMutex
	teams map[string]*domain.Team
}

func NewTeamRepository() *TeamRepository {
	return &TeamRepository{
		teams: make(map[string]*domain.Team),
	}
}

func (r *TeamRepository) Create(ctx context.Context, team domain.Team) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.teams[team.Name]; exists {
		return errors.New("team already exists")
	}

	// Создаём копию команды
	teamCopy := team
	r.teams[team.Name] = &teamCopy
	return nil
}

func (r *TeamRepository) GetByName(ctx context.Context, teamName string) (*domain.Team, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	team, exists := r.teams[teamName]
	if !exists {
		return nil, errors.New("team not found")
	}

	// Возвращаем копию
	teamCopy := *team
	membersCopy := make([]domain.TeamMember, len(team.Members))
	copy(membersCopy, team.Members)
	teamCopy.Members = membersCopy

	return &teamCopy, nil
}

// UpdateMember обновляет статус активности участника команды
func (r *TeamRepository) UpdateMember(ctx context.Context, teamName string, userID string, isActive bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	team, exists := r.teams[teamName]
	if !exists {
		return errors.New("team not found")
	}

	// Ищем участника в команде и обновляем его статус
	for i := range team.Members {
		if team.Members[i].UserID == userID {
			team.Members[i].IsActive = isActive
			return nil
		}
	}

	return errors.New("member not found in team")
}
