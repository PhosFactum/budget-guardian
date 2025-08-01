package main

import (
	"log"

	"github.com/PhosFactum/budget-guardian/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
