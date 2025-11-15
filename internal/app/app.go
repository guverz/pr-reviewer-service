package app

import (
	"context"
	"fmt"

	"github.com/guverz/pr-reviewer-service/internal/api"
	"github.com/guverz/pr-reviewer-service/internal/config"
	"github.com/guverz/pr-reviewer-service/internal/httpserver"
	"github.com/guverz/pr-reviewer-service/internal/repository/inmemory"
	"github.com/guverz/pr-reviewer-service/internal/service"
)

type Application struct {
	cfg    *config.Config
	server *httpserver.Server
}

func New(opts ...Option) (*Application, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	// Инициализируем репозитории
	repos := inmemory.NewRepositories()

	// Инициализируем сервисы
	reviewerSelector := service.NewReviewerSelector()
	teamService := service.NewTeamService(repos.Team, repos.User, repos.Transaction)
	userService := service.NewUserService(repos.User, repos.Team)
	pullRequestService := service.NewPullRequestService(repos.PullRequest, repos.User, reviewerSelector)

	// Создаём роутер
	router := api.NewRouter(teamService, userService, pullRequestService)

	// Инициализируем HTTP сервер
	server, err := httpserver.New(cfg, router)
	if err != nil {
		return nil, fmt.Errorf("init http server: %w", err)
	}

	app := &Application{
		cfg:    cfg,
		server: server,
	}

	for _, opt := range opts {
		opt(app)
	}

	return app, nil
}

func (a *Application) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return a.server.Serve(ctx)
}
