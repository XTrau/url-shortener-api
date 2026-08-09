package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"urlshortener/internal/apperrors"
	"urlshortener/internal/cache"
	"urlshortener/internal/database"
	"urlshortener/internal/domain"
	"urlshortener/internal/usecases"
)

type ShortenerRoutes struct {
	useCases usecases.UrlUseCases
}

func NewShortenerRoutes(urlRepo database.UrlRepository, urlCache cache.UrlCacher) *ShortenerRoutes {
	useCases := usecases.NewUrlUseCases(urlRepo, urlCache)
	return &ShortenerRoutes{useCases}
}

func (sr *ShortenerRoutes) RegisterRoutes(mux *http.ServeMux) {
	slog.Debug("Registering shortener routes")
	mux.HandleFunc("POST /short", sr.ShortenerHandler)
	mux.HandleFunc("GET /s/{slug}", sr.RedirectHandler)
}

func (sr *ShortenerRoutes) ShortenerHandler(w http.ResponseWriter, r *http.Request) {
	var urlReq domain.Url

	if err := json.NewDecoder(r.Body).Decode(&urlReq); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	slog.Debug("Request to Short url", slog.Any("Request Body", urlReq))

	slug, slugID, err := sr.useCases.GetSlug(urlReq.Url)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		panic(err)
	}

	urlResp := domain.Slug{
		Slug:   slug,
		SlugID: slugID,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(urlResp); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		panic(err)
	}
}

func (sr *ShortenerRoutes) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	slugID, err := strconv.Atoi(r.URL.Query().Get("i"))

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	slog.Debug("slug request", slog.String("slug", slug))

	url, err := sr.useCases.GetUrl(slug, slugID)

	if err != nil {
		if errors.Is(err, apperrors.ErrUrlNotFound) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			panic(err)
		}
		return
	}

	w.Header().Set("Location", url)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusMovedPermanently)
}
