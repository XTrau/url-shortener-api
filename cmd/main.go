package main

import (
	"log/slog"
	"urlshortener/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error(err.Error())
	}
}
