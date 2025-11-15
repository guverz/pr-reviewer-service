package service

import (
	"math/rand"
	"time"

	"github.com/guverz/pr-reviewer-service/internal/domain"
)

// ReviewerSelector выбирает ревьюеров из команды
type ReviewerSelector struct {
	rng *rand.Rand
}

func NewReviewerSelector() *ReviewerSelector {
	return &ReviewerSelector{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// SelectReviewers выбирает до 2 активных ревьюеров из команды, исключая автора
func (rs *ReviewerSelector) SelectReviewers(teamMembers []domain.User, excludeUserID string, maxCount int) []string {
	if maxCount <= 0 {
		maxCount = 2
	}

	// Фильтруем активных пользователей, исключая автора
	candidates := make([]domain.User, 0)
	for _, member := range teamMembers {
		if member.IsActive && member.ID != excludeUserID {
			candidates = append(candidates, member)
		}
	}

	if len(candidates) == 0 {
		return []string{}
	}

	// Выбираем случайных ревьюеров
	count := len(candidates)
	if count > maxCount {
		count = maxCount
	}

	// Перемешиваем кандидатов
	shuffled := make([]domain.User, len(candidates))
	copy(shuffled, candidates)
	rs.rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// Берём первых count
	reviewers := make([]string, 0, count)
	for i := 0; i < count; i++ {
		reviewers = append(reviewers, shuffled[i].ID)
	}

	return reviewers
}

