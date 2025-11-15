package repository

import (
	"context"

	"github.com/guverz/pr-reviewer-service/internal/domain"
)

type TeamRepository interface {
	Create(ctx context.Context, team domain.Team) error
	GetByName(ctx context.Context, teamName string) (*domain.Team, error)
	UpdateMember(ctx context.Context, teamName string, userID string, isActive bool) error
}

type UserRepository interface {
	UpsertTeamMembers(ctx context.Context, teamName string, members []domain.TeamMember) error
	GetByID(ctx context.Context, userID string) (*domain.User, error)
	SetActive(ctx context.Context, userID string, isActive bool) (*domain.User, error)
	ListByTeam(ctx context.Context, teamName string, onlyActive bool) ([]domain.User, error)
}

type PullRequestRepository interface {
	Create(ctx context.Context, pr domain.PullRequest) error
	GetByID(ctx context.Context, prID string) (*domain.PullRequest, error)
	Update(ctx context.Context, pr domain.PullRequest) error
	ListByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequest, error)
}

type TransactionManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
