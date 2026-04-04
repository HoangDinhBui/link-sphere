package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/linksphere/pkg/middleware"
	"github.com/linksphere/pkg/response"
	"github.com/linksphere/user-service/internal/model"
	"github.com/linksphere/user-service/internal/service"
)

// UserHandler handles HTTP requests for user operations.
type UserHandler struct {
	svc *service.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// Register handles user registration.
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.svc.Register(r.Context(), &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, "user registered successfully", user)
}

// Follow handles following a user.
func (h *UserHandler) Follow(w http.ResponseWriter, r *http.Request) {
	var req model.FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := middleware.GetUserID(r.Context())
	if err := h.svc.Follow(r.Context(), userID, req.TargetUserID); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, "followed successfully", nil)
}

// Unfollow handles unfollowing a user.
func (h *UserHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	var req model.FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := middleware.GetUserID(r.Context())
	if err := h.svc.Unfollow(r.Context(), userID, req.TargetUserID); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, "unfollowed successfully", nil)
}

// GetProfile returns the current user's profile.
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	user, err := h.svc.GetByID(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "user not found")
		return
	}
	response.JSON(w, http.StatusOK, user)
}

// GetUserByID returns a user by their ID.
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "user not found")
		return
	}
	response.JSON(w, http.StatusOK, user)
}

// GetFollowing returns the list of IDs that the user follows.
func (h *UserHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	following, err := h.svc.GetFollowing(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, following)
}
