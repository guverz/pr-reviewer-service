package main

import (
	"log"

	"github.com/guverz/pr-reviewer-service/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("application stopped with error: %v", err)
	}
}



