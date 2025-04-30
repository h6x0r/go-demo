package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/h6x0r/go-demo/internal/database"
	"github.com/h6x0r/go-demo/internal/handlers"
	"github.com/h6x0r/go-demo/internal/models"
	"github.com/h6x0r/go-demo/internal/repository"
)

func TestMain(m *testing.M) {
	database.Connect()
	database.Migrate()

	code := m.Run()

	database.DB.Close()
	os.Exit(code)
}

// Для очистки таблицы между тестами
func truncateUsersTable() {
	_, _ = database.DB.Exec("TRUNCATE TABLE users")
}

// TestUserIntegration_Create проверяет создание пользователя
func TestUserIntegration_Create(t *testing.T) {
	truncateUsersTable()

	gin.SetMode(gin.TestMode)
	repo := repository.NewUserRepository(database.DB)
	userHandler := handlers.NewUserHandler(repo)

	r := gin.Default()
	r.POST("/users", userHandler.CreateUser)

	body := []byte(`{
		"firstname": "Bob",
		"lastname": "Marley",
		"email": "bob@example.com",
		"age": 35
	}`)

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Ожидается статус 201 при успешном создании пользователя")

	var created models.User
	err := json.Unmarshal(w.Body.Bytes(), &created)
	assert.Nil(t, err)
	assert.NotEqual(t, uuid.Nil, created.ID)
	assert.NotEmpty(t, created.ID, "ID пользователя не должен быть пустым")
}

// TestUserIntegration_Get проверяет получение пользователя
func TestUserIntegration_Get(t *testing.T) {
	truncateUsersTable()

	gin.SetMode(gin.TestMode)
	repo := repository.NewUserRepository(database.DB)
	userHandler := handlers.NewUserHandler(repo)
	r := gin.Default()

	r.POST("/users", userHandler.CreateUser)
	r.GET("/user/:id", userHandler.GetUser)

	body := []byte(`{
		"firstname": "Alice",
		"lastname": "Wonder",
		"email": "alice@example.com",
		"age": 22
	}`)

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ws := httptest.NewRecorder()
	r.ServeHTTP(ws, req)
	assert.Equal(t, http.StatusCreated, ws.Code)

	var created models.User
	_ = json.Unmarshal(ws.Body.Bytes(), &created)
	userID := created.ID

	reqGet, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user/%s", userID.String()), nil)
	wg := httptest.NewRecorder()
	r.ServeHTTP(wg, reqGet)

	assert.Equal(t, http.StatusOK, wg.Code, "Ожидается статус 200 при успешном получении пользователя")

	var gotUser models.User
	_ = json.Unmarshal(wg.Body.Bytes(), &gotUser)
	assert.Equal(t, userID, gotUser.ID, "ID полученного пользователя не совпадает с созданным")
}

// TestUserIntegration_Update проверяет обновление пользователя
func TestUserIntegration_Update(t *testing.T) {
	truncateUsersTable()

	gin.SetMode(gin.TestMode)
	repo := repository.NewUserRepository(database.DB)
	userHandler := handlers.NewUserHandler(repo)
	r := gin.Default()

	r.POST("/users", userHandler.CreateUser)
	r.GET("/user/:id", userHandler.GetUser)
	r.PATCH("/user/:id", userHandler.UpdateUser)

	bodyCreate := []byte(`{
		"firstname": "Carl",
		"lastname": "Johnson",
		"email": "carl@example.com",
		"age": 28
	}`)

	reqCreate, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(bodyCreate))
	reqCreate.Header.Set("Content-Type", "application/json")
	wc := httptest.NewRecorder()
	r.ServeHTTP(wc, reqCreate)
	assert.Equal(t, http.StatusCreated, wc.Code)

	var created models.User
	_ = json.Unmarshal(wc.Body.Bytes(), &created)
	userID := created.ID

	updateBody := []byte(`{
		"firstname": "Carlos",
		"age": 29
	}`)

	reqUpdate, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/user/%s", userID.String()), bytes.NewBuffer(updateBody))
	reqUpdate.Header.Set("Content-Type", "application/json")
	wu := httptest.NewRecorder()
	r.ServeHTTP(wu, reqUpdate)

	assert.Equal(t, http.StatusOK, wu.Code, "Ожидается статус 200 при успешном обновлении")

	reqGet, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user/%s", userID.String()), nil)
	wg := httptest.NewRecorder()
	r.ServeHTTP(wg, reqGet)
	assert.Equal(t, http.StatusOK, wg.Code)

	var updated models.UpdateUser
	_ = json.Unmarshal(wg.Body.Bytes(), &updated)
	assert.NotNil(t, updated.Firstname, "Firstname не должно быть nil")
	assert.Equal(t, "Carlos", *updated.Firstname, "Имя не обновилось")

	var age uint = 29
	assert.NotNil(t, updated.Age, "Age не должно быть nil")
	assert.Equal(t, &age, updated.Age, "Возраст не обновился")
}

// TestUserIntegration_Delete проверяет удаление пользователя (если метод реализован)
func TestUserIntegration_Delete(t *testing.T) {
	truncateUsersTable()

	gin.SetMode(gin.TestMode)
	repo := repository.NewUserRepository(database.DB)
	userHandler := handlers.NewUserHandler(repo)
	r := gin.Default()

	r.POST("/users", userHandler.CreateUser)
	r.GET("/user/:id", userHandler.GetUser)
	r.DELETE("/user/:id", userHandler.DeleteUser)

	bodyCreate := []byte(`{
		"firstname": "Diana",
		"lastname": "Prince",
		"email": "diana@example.com",
		"age": 25
	}`)

	reqCreate, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(bodyCreate))
	reqCreate.Header.Set("Content-Type", "application/json")
	wc := httptest.NewRecorder()
	r.ServeHTTP(wc, reqCreate)
	assert.Equal(t, http.StatusCreated, wc.Code)

	var created models.User
	_ = json.Unmarshal(wc.Body.Bytes(), &created)
	userID := created.ID

	reqDelete, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/user/%s", userID.String()), nil)
	wd := httptest.NewRecorder()
	r.ServeHTTP(wd, reqDelete)

	assert.True(t, wd.Code == http.StatusNoContent || wd.Code == http.StatusOK, "Ожидался статус 204 или 200")

	reqGet, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/user/%s", userID.String()), nil)
	wg := httptest.NewRecorder()
	r.ServeHTTP(wg, reqGet)

	assert.Equal(t, http.StatusNotFound, wg.Code, "После удаления пользователь должен отсутствовать")
}
