package service

import (
	"context"
	"time"

	"github.com/guverz/pr-reviewer-service/internal/domain"
	"github.com/guverz/pr-reviewer-service/internal/repository"
)

type PullRequestService struct {
	prRepo           repository.PullRequestRepository
	userRepo         repository.UserRepository
	reviewerSelector *ReviewerSelector
}

func NewPullRequestService(
	prRepo repository.PullRequestRepository,
	userRepo repository.UserRepository,
	reviewerSelector *ReviewerSelector,
) *PullRequestService {
	return &PullRequestService{
		prRepo:           prRepo,
		userRepo:         userRepo,
		reviewerSelector: reviewerSelector,
	}
}

// CreatePR создаёт PR и автоматически назначает до 2 ревьюеров из команды автора
func (s *PullRequestService) CreatePR(ctx context.Context, prID, prName, authorID string) (*domain.PullRequest, error) {
	// Проверяем, существует ли PR
	existing, err := s.prRepo.GetByID(ctx, prID)
	if err == nil && existing != nil {
		return nil, domain.NewDomainError(domain.ErrorCodePRExists, "PR id already exists")
	}

	// Получаем автора
	author, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil {
		return nil, domain.NewDomainError(domain.ErrorCodeNotFound, "author not found")
	}

	// Получаем активных участников команды автора
	teamMembers, err := s.userRepo.ListByTeam(ctx, author.TeamName, true)
	if err != nil {
		return nil, domain.NewDomainError(domain.ErrorCodeNotFound, "team not found")
	}

	// Выбираем ревьюеров
	reviewers := s.reviewerSelector.SelectReviewers(teamMembers, authorID, 2)

	// Создаём PR
	pr := domain.PullRequest{
		ID:                prID,
		Name:              prName,
		AuthorID:          authorID,
		Status:            domain.PullRequestStatusOpen,
		AssignedReviewers: reviewers,
		NeedMoreReviewers: len(reviewers) < 2,
		CreatedAt:         time.Now(),
		MergedAt:          nil,
	}

	if err := s.prRepo.Create(ctx, pr); err != nil {
		return nil, err
	}

	return &pr, nil
}

// MergePR помечает PR как MERGED (идемпотентная операция)
func (s *PullRequestService) MergePR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, domain.NewDomainError(domain.ErrorCodeNotFound, "PR not found")
	}

	// Если уже merged, просто возвращаем текущее состояние (идемпотентность)
	if pr.IsMerged() {
		return pr, nil
	}

	// Помечаем как merged
	now := time.Now()
	pr.Status = domain.PullRequestStatusMerged
	pr.MergedAt = &now

	if err := s.prRepo.Update(ctx, *pr); err != nil {
		return nil, err
	}

	return pr, nil
}

// ReassignReviewer переназначает ревьюера на другого из его команды
func (s *PullRequestService) ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*domain.PullRequest, string, error) {
	// Получаем PR
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, "", domain.NewDomainError(domain.ErrorCodeNotFound, "PR not found")
	}

	// Проверяем, что PR не merged
	if pr.IsMerged() {
		return nil, "", domain.NewDomainError(domain.ErrorCodePRMerged, "cannot reassign on merged PR")
	}

	// Проверяем, что старый ревьюер назначен
	if !pr.HasReviewer(oldReviewerID) {
		return nil, "", domain.NewDomainError(domain.ErrorCodeNotAssigned, "reviewer is not assigned to this PR")
	}

	// Получаем старого ревьюера для определения его команды
	oldReviewer, err := s.userRepo.GetByID(ctx, oldReviewerID)
	if err != nil {
		return nil, "", domain.NewDomainError(domain.ErrorCodeNotFound, "reviewer not found")
	}

	// Получаем активных участников команды старого ревьюера
	teamMembers, err := s.userRepo.ListByTeam(ctx, oldReviewer.TeamName, true)
	if err != nil {
		return nil, "", domain.NewDomainError(domain.ErrorCodeNotFound, "team not found")
	}

	// Исключаем только автора и заменяемого ревьюера
	// Уже назначенные ревьюеры могут быть выбраны снова (по ТЗ это допустимо)
	excludeIDs := make(map[string]bool)
	excludeIDs[oldReviewerID] = true
	excludeIDs[pr.AuthorID] = true

	// Фильтруем кандидатов (только активные, исключая автора и заменяемого)
	candidates := make([]domain.User, 0)
	for _, member := range teamMembers {
		if member.IsActive && !excludeIDs[member.ID] {
			candidates = append(candidates, member)
		}
	}

	if len(candidates) == 0 {
		return nil, "", domain.NewDomainError(domain.ErrorCodeNoCandidate, "no active replacement candidate in team")
	}

	// Выбираем случайного кандидата
	selected := s.reviewerSelector.SelectReviewers(candidates, "", 1)
	if len(selected) == 0 {
		return nil, "", domain.NewDomainError(domain.ErrorCodeNoCandidate, "no active replacement candidate in team")
	}

	newReviewerID := selected[0]

	// Заменяем ревьюера
	pr.ReplaceReviewer(oldReviewerID, newReviewerID)
	pr.NeedMoreReviewers = len(pr.AssignedReviewers) < 2

	if err := s.prRepo.Update(ctx, *pr); err != nil {
		return nil, "", err
	}

	return pr, newReviewerID, nil
}

// GetPRsByReviewer получает список PR, где пользователь назначен ревьюером
func (s *PullRequestService) GetPRsByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequest, error) {
	// Проверяем, что пользователь существует
	_, err := s.userRepo.GetByID(ctx, reviewerID)
	if err != nil {
		return nil, domain.NewDomainError(domain.ErrorCodeNotFound, "user not found")
	}

	prs, err := s.prRepo.ListByReviewer(ctx, reviewerID)
	if err != nil {
		return nil, err
	}

	return prs, nil
}
