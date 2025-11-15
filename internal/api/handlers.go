package api

import (
	"encoding/json"
	"net/http"

	"github.com/guverz/pr-reviewer-service/internal/domain"
	"github.com/guverz/pr-reviewer-service/internal/service"
)

type Handlers struct {
	teamService        *service.TeamService
	userService        *service.UserService
	pullRequestService *service.PullRequestService
}

func NewHandlers(
	teamService *service.TeamService,
	userService *service.UserService,
	pullRequestService *service.PullRequestService,
) *Handlers {
	return &Handlers{
		teamService:        teamService,
		userService:        userService,
		pullRequestService: pullRequestService,
	}
}

// POST /team/add
func (h *Handlers) AddTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TeamDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, domain.NewDomainError(domain.ErrorCodeNotFound, "invalid request body"))
		return
	}

	team := ToTeam(req)
	createdTeam, err := h.teamService.CreateTeam(r.Context(), team)
	if err != nil {
		WriteError(w, err)
		return
	}

	response := TeamResponse{
		Team: ToTeamDTO(*createdTeam),
	}
	WriteJSON(w, http.StatusCreated, response)
}

// GET /team/get
func (h *Handlers) GetTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		WriteError(w, domain.NewDomainError(domain.ErrorCodeNotFound, "team_name parameter is required"))
		return
	}

	team, err := h.teamService.GetTeam(r.Context(), teamName)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, ToTeamDTO(*team))
}

// POST /users/setIsActive
func (h *Handlers) SetUserActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SetActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, domain.NewDomainError(domain.ErrorCodeNotFound, "invalid request body"))
		return
	}

	user, err := h.userService.SetActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		WriteError(w, err)
		return
	}

	response := UserResponse{
		User: ToUserDTO(*user),
	}
	WriteJSON(w, http.StatusOK, response)
}

// GET /users/getReview
func (h *Handlers) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		WriteError(w, domain.NewDomainError(domain.ErrorCodeNotFound, "user_id parameter is required"))
		return
	}

	prs, err := h.pullRequestService.GetPRsByReviewer(r.Context(), userID)
	if err != nil {
		WriteError(w, err)
		return
	}

	prDTOs := make([]PullRequestShortDTO, len(prs))
	for i, pr := range prs {
		prDTOs[i] = ToPullRequestShortDTO(pr)
	}

	response := GetReviewResponse{
		UserID:       userID,
		PullRequests: prDTOs,
	}
	WriteJSON(w, http.StatusOK, response)
}

// POST /pullRequest/create
func (h *Handlers) CreatePR(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, domain.NewDomainError(domain.ErrorCodeNotFound, "invalid request body"))
		return
	}

	pr, err := h.pullRequestService.CreatePR(r.Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		WriteError(w, err)
		return
	}

	response := PullRequestResponse{
		PR: ToPullRequestDTO(*pr),
	}
	WriteJSON(w, http.StatusCreated, response)
}

// POST /pullRequest/merge
func (h *Handlers) MergePR(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, domain.NewDomainError(domain.ErrorCodeNotFound, "invalid request body"))
		return
	}

	pr, err := h.pullRequestService.MergePR(r.Context(), req.PullRequestID)
	if err != nil {
		WriteError(w, err)
		return
	}

	response := PullRequestResponse{
		PR: ToPullRequestDTO(*pr),
	}
	WriteJSON(w, http.StatusOK, response)
}

// POST /pullRequest/reassign
func (h *Handlers) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, domain.NewDomainError(domain.ErrorCodeNotFound, "invalid request body"))
		return
	}

	pr, newReviewerID, err := h.pullRequestService.ReassignReviewer(r.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		WriteError(w, err)
		return
	}

	response := ReassignResponse{
		PR:         ToPullRequestDTO(*pr),
		ReplacedBy: newReviewerID,
	}
	WriteJSON(w, http.StatusOK, response)
}
