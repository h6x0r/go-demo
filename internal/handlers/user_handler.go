package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/h6x0r/go-demo/internal/models"
	"github.com/h6x0r/go-demo/internal/repository"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

// CreateUser POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = uuid.New()
	user.Created = time.Now().UTC()

	if user.Firstname == "" || user.Lastname == "" || user.Email == "" || user.Age == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "firstname, lastname, email и age обязательны"})
		return
	}

	if err := h.repo.CreateUser(context.Background(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser GET /user/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат UUID"})
		return
	}

	user, err := h.repo.GetUserByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser PATCH /user/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат UUID"})
		return
	}

	var userUpdate models.UpdateUser
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.repo.GetUserByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	if userUpdate.Firstname != nil && *userUpdate.Firstname != "" {
		user.Firstname = *userUpdate.Firstname
	}

	if userUpdate.Lastname != nil && *userUpdate.Lastname != "" {
		user.Lastname = *userUpdate.Lastname
	}

	if userUpdate.Email != nil && *userUpdate.Email != "" {
		user.Email = *userUpdate.Email
	}

	if userUpdate.Age != nil && *userUpdate.Age != 0 {
		user.Age = *userUpdate.Age
	}

	err = h.repo.UpdateUser(context.Background(), *user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser DELETE /user/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат UUID"})
		return
	}

	_, err = h.repo.GetUserByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	err = h.repo.DeleteUser(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
