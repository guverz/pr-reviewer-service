package api

import (
	"net/http"

	"github.com/guverz/pr-reviewer-service/internal/service"
)

func NewRouter(
	teamService *service.TeamService,
	userService *service.UserService,
	pullRequestService *service.PullRequestService,
) http.Handler {
	mux := http.NewServeMux()

	handlers := NewHandlers(teamService, userService, pullRequestService)

	// Teams endpoints
	mux.HandleFunc("/team/add", handlers.AddTeam)
	mux.HandleFunc("/team/get", handlers.GetTeam)

	// Users endpoints
	mux.HandleFunc("/users/setIsActive", handlers.SetUserActive)
	mux.HandleFunc("/users/getReview", handlers.GetUserReviews)

	// PullRequests endpoints
	mux.HandleFunc("/pullRequest/create", handlers.CreatePR)
	mux.HandleFunc("/pullRequest/merge", handlers.MergePR)
	mux.HandleFunc("/pullRequest/reassign", handlers.ReassignReviewer)
	
	// Health check
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	return mux
}
