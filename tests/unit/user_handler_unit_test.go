package unit

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/h6x0r/go-demo/internal/handlers"
	"github.com/h6x0r/go-demo/internal/models"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(ctx context.Context, user models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepo) UpdateUser(ctx context.Context, user models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockUserRepo)
	handler := handlers.NewUserHandler(mockRepo)

	r := gin.Default()
	r.POST("/users", handler.CreateUser)

	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("models.User")).
		Return(nil)

	requestBody := []byte(`{
        "firstname": "John",
        "lastname": "Doe",
        "email": "john.doe@example.com",
        "age": 25
    }`)

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	mockRepo.AssertNumberOfCalls(t, "CreateUser", 1)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockUserRepo)
	handler := handlers.NewUserHandler(mockRepo)

	r := gin.Default()
	r.POST("/users", handler.CreateUser)

	requestBody := []byte(`{
        "firstname": "",
        "lastname": "Doe",
        "email": "not_an_email",
        "age": 25
    }`)

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRepo.AssertNotCalled(t, "CreateUser")
}

func TestGetUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockUserRepo)
	handler := handlers.NewUserHandler(mockRepo)

	r := gin.Default()
	r.GET("/user/:id", handler.GetUser)

	userID := uuid.New()
	user := &models.User{
		ID:        userID,
		Firstname: "Jane",
		Lastname:  "Doe",
		Email:     "jane@example.com",
		Age:       32,
		Created:   time.Now(),
	}

	mockRepo.On("GetUserByID", mock.Anything, userID).
		Return(user, nil)

	req, _ := http.NewRequest(http.MethodGet, "/user/"+userID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertNumberOfCalls(t, "GetUserByID", 1)
	mockRepo.AssertExpectations(t)
}

func TestGetUser_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockUserRepo)
	handler := handlers.NewUserHandler(mockRepo)

	r := gin.Default()
	r.GET("/user/:id", handler.GetUser)

	userID := uuid.New()
	mockRepo.On("GetUserByID", mock.Anything, userID).
		Return(nil, errors.New("Not found"))

	req, _ := http.NewRequest(http.MethodGet, "/user/"+userID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockUserRepo)
	handler := handlers.NewUserHandler(mockRepo)

	r := gin.Default()
	r.PATCH("/user/:id", handler.UpdateUser)

	userID := uuid.New()
	existingUser := &models.User{
		ID:        userID,
		Firstname: "Mike",
		Lastname:  "Scott",
		Email:     "mike.scott@example.com",
		Age:       40,
		Created:   time.Now(),
	}

	mockRepo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil)
	mockRepo.On("UpdateUser", mock.Anything, mock.AnythingOfType("models.User")).
		Return(nil)

	requestBody := []byte(`{
        "firstname": "Michael",
        "age": 42
    }`)

	req, _ := http.NewRequest(http.MethodPatch, "/user/"+userID.String(), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertNumberOfCalls(t, "GetUserByID", 1)
	mockRepo.AssertNumberOfCalls(t, "UpdateUser", 1)
	mockRepo.AssertExpectations(t)
}
