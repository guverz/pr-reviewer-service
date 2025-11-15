package api

import (
	"time"

	"github.com/guverz/pr-reviewer-service/internal/domain"
)

// Team DTO
type TeamMemberDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type TeamDTO struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}

type TeamResponse struct {
	Team TeamDTO `json:"team"`
}

// User DTO
type UserDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type UserResponse struct {
	User UserDTO `json:"user"`
}

type SetActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

// PullRequest DTO
type PullRequestDTO struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	CreatedAt         *string  `json:"createdAt,omitempty"`
	MergedAt          *string  `json:"mergedAt,omitempty"`
}

type PullRequestShortDTO struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

type PullRequestResponse struct {
	PR PullRequestDTO `json:"pr"`
}

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type ReassignRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

type ReassignResponse struct {
	PR         PullRequestDTO `json:"pr"`
	ReplacedBy string         `json:"replaced_by"`
}

type GetReviewResponse struct {
	UserID       string                `json:"user_id"`
	PullRequests []PullRequestShortDTO `json:"pull_requests"`
}

// Error DTO
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// Конвертеры из domain в DTO
func ToTeamMemberDTO(m domain.TeamMember) TeamMemberDTO {
	return TeamMemberDTO{
		UserID:   m.UserID,
		Username: m.Username,
		IsActive: m.IsActive,
	}
}

func ToTeamDTO(t domain.Team) TeamDTO {
	members := make([]TeamMemberDTO, len(t.Members))
	for i, m := range t.Members {
		members[i] = ToTeamMemberDTO(m)
	}
	return TeamDTO{
		TeamName: t.Name,
		Members:  members,
	}
}

func ToUserDTO(u domain.User) UserDTO {
	return UserDTO{
		UserID:   u.ID,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func ToPullRequestDTO(pr domain.PullRequest) PullRequestDTO {
	createdAt := formatTime(pr.CreatedAt)
	var mergedAt *string
	if pr.MergedAt != nil {
		formatted := formatTime(*pr.MergedAt)
		mergedAt = &formatted
	}

	return PullRequestDTO{
		PullRequestID:     pr.ID,
		PullRequestName:   pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         &createdAt,
		MergedAt:          mergedAt,
	}
}

func ToPullRequestShortDTO(pr domain.PullRequest) PullRequestShortDTO {
	return PullRequestShortDTO{
		PullRequestID:   pr.ID,
		PullRequestName: pr.Name,
		AuthorID:        pr.AuthorID,
		Status:          string(pr.Status),
	}
}

// Конвертеры из DTO в domain
func ToTeam(dto TeamDTO) domain.Team {
	members := make([]domain.TeamMember, len(dto.Members))
	for i, m := range dto.Members {
		members[i] = domain.TeamMember{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}
	return domain.Team{
		Name:    dto.TeamName,
		Members: members,
	}
}



