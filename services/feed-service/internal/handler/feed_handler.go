package handler

import (
	"encoding/json"
	"net/http"

	"github.com/linksphere/feed-service/internal/service"
	"github.com/linksphere/pkg/middleware"
	"github.com/linksphere/pkg/response"
)

// FeedRequest represents the feed request body.
type FeedRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// FeedHandler handles HTTP requests for feed operations.
type FeedHandler struct {
	svc *service.FeedService
}

// NewFeedHandler creates a new FeedHandler.
func NewFeedHandler(svc *service.FeedService) *FeedHandler {
	return &FeedHandler{svc: svc}
}

// GetFeed returns the news feed for the authenticated user.
func (h *FeedHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	var req FeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 50 {
		req.Limit = 10
	}

	userID := middleware.GetUserID(r.Context())
	feed, err := h.svc.GetFeed(r.Context(), userID, req.Page, req.Limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, feed)
}
