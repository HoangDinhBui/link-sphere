package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/linksphere/pkg/middleware"
	"github.com/linksphere/pkg/response"
	"github.com/linksphere/post-service/internal/model"
	"github.com/linksphere/post-service/internal/service"
)

// PostHandler handles HTTP requests for post operations.
type PostHandler struct {
	svc *service.PostService
}

// NewPostHandler creates a new PostHandler.
func NewPostHandler(svc *service.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

// Create handles creating a new post.
func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := middleware.GetUserID(r.Context())
	post, err := h.svc.Create(r.Context(), userID, &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, "post created successfully", post)
}

// ListRequest represents the list posts request body.
type ListRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// DetailRequest represents the get post detail request body.
type DetailRequest struct {
	PostID string `json:"postId"`
}

// List returns a paginated list of posts.
func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	var req ListRequest
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

	posts, err := h.svc.List(r.Context(), req.Page, req.Limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, posts)
}

// GetByID returns a single post by ID.
func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	var req DetailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	post, err := h.svc.GetByID(r.Context(), req.PostID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "post not found")
		return
	}
	response.JSON(w, http.StatusOK, post)
}

// Like handles liking a post.
func (h *PostHandler) Like(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")
	userID := middleware.GetUserID(r.Context())

	if err := h.svc.Like(r.Context(), postID, userID); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, "post liked successfully", nil)
}

// GetByUserIDsRequest represents the request body to retrieve posts by user IDs.
type GetByUserIDsRequest struct {
	UserIDs []string `json:"userIds"`
	Page    int      `json:"page"`
	Limit   int      `json:"limit"`
}

// GetByUserIDs returns posts for multiple user IDs.
func (h *PostHandler) GetByUserIDs(w http.ResponseWriter, r *http.Request) {
	var req GetByUserIDsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}

	posts, err := h.svc.GetByUserIDs(r.Context(), req.UserIDs, req.Page, req.Limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, posts)
}
