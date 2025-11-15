package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/guverz/pr-reviewer-service/internal/domain"
)

type PullRequestRepository struct {
	mu  sync.RWMutex
	prs map[string]*domain.PullRequest
}

func NewPullRequestRepository() *PullRequestRepository {
	return &PullRequestRepository{
		prs: make(map[string]*domain.PullRequest),
	}
}

func (r *PullRequestRepository) Create(ctx context.Context, pr domain.PullRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.prs[pr.ID]; exists {
		return errors.New("PR already exists")
	}

	// Создаём копию PR
	prCopy := pr
	reviewersCopy := make([]string, len(pr.AssignedReviewers))
	copy(reviewersCopy, pr.AssignedReviewers)
	prCopy.AssignedReviewers = reviewersCopy

	r.prs[pr.ID] = &prCopy
	return nil
}

func (r *PullRequestRepository) GetByID(ctx context.Context, prID string) (*domain.PullRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pr, exists := r.prs[prID]
	if !exists {
		return nil, errors.New("PR not found")
	}

	// Возвращаем копию
	prCopy := *pr
	reviewersCopy := make([]string, len(pr.AssignedReviewers))
	copy(reviewersCopy, pr.AssignedReviewers)
	prCopy.AssignedReviewers = reviewersCopy

	return &prCopy, nil
}

func (r *PullRequestRepository) Update(ctx context.Context, pr domain.PullRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.prs[pr.ID]; !exists {
		return errors.New("PR not found")
	}

	// Обновляем PR
	prCopy := pr
	reviewersCopy := make([]string, len(pr.AssignedReviewers))
	copy(reviewersCopy, pr.AssignedReviewers)
	prCopy.AssignedReviewers = reviewersCopy

	r.prs[pr.ID] = &prCopy
	return nil
}

func (r *PullRequestRepository) ListByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	prs := make([]domain.PullRequest, 0)
	for _, pr := range r.prs {
		if pr.HasReviewer(reviewerID) {
			prCopy := *pr
			reviewersCopy := make([]string, len(pr.AssignedReviewers))
			copy(reviewersCopy, pr.AssignedReviewers)
			prCopy.AssignedReviewers = reviewersCopy
			prs = append(prs, prCopy)
		}
	}

	return prs, nil
}



