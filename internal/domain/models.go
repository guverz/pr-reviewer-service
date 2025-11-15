package domain

import "time"

type TeamMember struct {
	UserID   string
	Username string
	IsActive bool
}

type Team struct {
	Name    string
	Members []TeamMember
}

type User struct {
	ID       string
	Username string
	TeamName string
	IsActive bool
}

type PullRequest struct {
	ID                string
	Name              string
	AuthorID          string
	Status            PullRequestStatus
	AssignedReviewers []string
	NeedMoreReviewers bool
	CreatedAt         time.Time
	MergedAt          *time.Time
}

func (pr *PullRequest) HasReviewer(userID string) bool {
	for _, reviewerID := range pr.AssignedReviewers {
		if reviewerID == userID {
			return true
		}
	}
	return false
}

func (pr *PullRequest) ReplaceReviewer(oldReviewer, newReviewer string) bool {
	for idx, reviewerID := range pr.AssignedReviewers {
		if reviewerID == oldReviewer {
			pr.AssignedReviewers[idx] = newReviewer
			return true
		}
	}
	return false
}

func (pr *PullRequest) IsMerged() bool {
	return pr.Status == PullRequestStatusMerged
}
