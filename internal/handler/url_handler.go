package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/1sh-repalto/url-monitoring-api/internal/middleware"
	"github.com/1sh-repalto/url-monitoring-api/internal/service"
	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(s *service.URLService) *URLHandler {
	return &URLHandler{service: s}
}

type RegisterURLRequest struct {
	URL string `json:"url"`
}

func (h *URLHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	parsedURL, err := url.ParseRequestURI(req.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	normalizedURL := parsedURL.Scheme + "://" + parsedURL.Host + parsedURL.Path
	if parsedURL.RawQuery != "" {
		normalizedURL += "?" + parsedURL.RawQuery
	}
	if parsedURL.Fragment != "" {
		normalizedURL += "#" + parsedURL.Fragment
	}
	normalizedURL = strings.TrimSuffix(normalizedURL, "/")

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.service.RegisterURL(normalizedURL, userID)
	if err != nil {
		if err.Error() == "URL already registered" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Failed to register URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "URL registered",
	})
}

func (h *URLHandler) GetURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	urls, err := h.service.GetURLByUser(userID)
	if err != nil {
		http.Error(w, "Failed to get URLs", http.StatusInternalServerError)
		return
	}
	if len(urls) == 0 {
		json.NewEncoder(w).Encode([]string{})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(urls); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *URLHandler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "urlID")
	if id == "" {
		http.Error(w, "Missing url ID", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteURL(id, userID)
	if err != nil {
		if err.Error() == "unauthorized: not the owner of the URL" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to delete url registry", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *URLHandler) GetURLLogs(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	urlID := chi.URLParam(r, "urlID")
	if urlID == "" {
		http.Error(w, "Missing URL ID", http.StatusNotFound)
		return
	}

	url, err := h.service.GetURLByID(urlID)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	if url.UserID != userID {
		http.Error(w, "Forbidded", http.StatusForbidden)
		return
	}

	logs, err := h.service.GetLogsByURLID(urlID)
	if err != nil {
		http.Error(w, "Failed to fetch logs", http.StatusInternalServerError)
		return
	}
	if len(logs) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
