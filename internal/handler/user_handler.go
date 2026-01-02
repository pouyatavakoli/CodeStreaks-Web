package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pouyatavakoli/CodeStreaks-web/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type AddUserRequest struct {
	CodeforcesHandle string `json:"codeforces_handle" binding:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    any `json:"data,omitempty"`
}

type LeaderboardResponse struct {
	Users      any `json:"users"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// AddUser godoc
// @Summary Add a new user
// @Description Add a user by their Codeforces handle
// @Tags users
// @Accept json
// @Produce json
// @Param user body AddUserRequest true "User handle"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users [post]
func (h *UserHandler) AddUser(c *gin.Context) {
	var req AddUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.userService.AddUser(req.CodeforcesHandle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "User added successfully",
		Data:    user,
	})
}

// GetLeaderboard godoc
// @Summary Get leaderboard
// @Description Get paginated leaderboard of users sorted by streak
// @Tags leaderboard
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(50)
// @Success 200 {object} LeaderboardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/leaderboard [get]
func (h *UserHandler) GetLeaderboard(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	users, total, err := h.userService.GetLeaderboard(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, LeaderboardResponse{
		Users:      users,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetUserByHandle godoc
// @Summary Get user by handle
// @Description Get a specific user by their Codeforces handle
// @Tags users
// @Produce json
// @Param handle path string true "Codeforces handle"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/{handle} [get]
func (h *UserHandler) GetUserByHandle(c *gin.Context) {
	handle := c.Param("handle")

	user, err := h.userService.GetUserByHandle(handle)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: user,
	})
}
