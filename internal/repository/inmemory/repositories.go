package inmemory

import (
	"github.com/guverz/pr-reviewer-service/internal/repository"
)

type Repositories struct {
	Team        repository.TeamRepository
	User        repository.UserRepository
	PullRequest repository.PullRequestRepository
	Transaction repository.TransactionManager
}

func NewRepositories() *Repositories {
	teamRepo := NewTeamRepository()
	userRepo := NewUserRepository()
	prRepo := NewPullRequestRepository()
	txMgr := NewTransactionManager()

	return &Repositories{
		Team:        teamRepo,
		User:        userRepo,
		PullRequest: prRepo,
		Transaction: txMgr,
	}
}

