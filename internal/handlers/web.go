package handlers

import (
	"log/slog"
	"net/http"
)

type WebHandlers struct {
}

func NewWebHandlers() *WebHandlers {
	return &WebHandlers{}
}

func (wh *WebHandlers) RegisterRoutes(mux *http.ServeMux) {
	slog.Debug("Registering web routes")

	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/", fs)
}
