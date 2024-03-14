package testutil_test

import (
	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/rest/v1"
	"github.com/dannyh79/whostodo/internal/sessions"
	"github.com/dannyh79/whostodo/internal/tasks"
	"github.com/gin-gonic/gin"
)

type MockTestSuite struct {
	Engine      *gin.Engine
	TaskRepo    *MockTaskRepository
	SessionRepo *MockSessionsRepository
}

func NewTestSuite() *MockTestSuite {
	gin.SetMode(gin.TestMode)
	engine := gin.Default()

	taskRepo := &MockTaskRepository{
		Data: make(map[int]repository.TaskSchema),
	}
	sessionRepo := &MockSessionsRepository{
		Data: make(map[string]Session),
	}
	tasksUsecase := tasks.InitTasksUsecase(taskRepo)
	sessionsUsecase := sessions.InitSessionsUsecase(sessionRepo)

	routes.AddRoutes(engine, tasksUsecase, sessionsUsecase)

	return &MockTestSuite{
		Engine:      engine,
		TaskRepo:    taskRepo,
		SessionRepo: sessionRepo,
	}
}
