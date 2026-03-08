package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

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

// List returns a paginated list of posts.
func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}

	posts, err := h.svc.List(r.Context(), page, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, posts)
}

// GetByID returns a single post by ID.
func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	post, err := h.svc.GetByID(r.Context(), id)
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
