package handler

import (
	"encoding/json"
	"net/http"

	"github.com/linksphere/comment-service/internal/model"
	"github.com/linksphere/comment-service/internal/service"
	"github.com/linksphere/pkg/middleware"
	"github.com/linksphere/pkg/response"
)

// CommentHandler handles HTTP requests for comment operations.
type CommentHandler struct {
	svc *service.CommentService
}

// NewCommentHandler creates a new CommentHandler.
func NewCommentHandler(svc *service.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

// Create handles creating a new comment.
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := middleware.GetUserID(r.Context())
	comment, err := h.svc.Create(r.Context(), userID, &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, "comment created successfully", comment)
}

// List returns comments for a post.
func (h *CommentHandler) List(w http.ResponseWriter, r *http.Request) {
	var req model.ListCommentsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 50 {
		req.Limit = 20
	}

	comments, err := h.svc.ListByPostID(r.Context(), req.PostID, req.Page, req.Limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, comments)
}
