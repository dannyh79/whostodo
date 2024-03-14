package main

import (
	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/rest/v1"
	"github.com/dannyh79/whostodo/internal/sessions"
	"github.com/dannyh79/whostodo/internal/tasks"
	"github.com/gin-gonic/gin"
)

func main() {
	taskRepo := repository.InitInMemoryTaskRepository()
	tasksUsecase := tasks.InitTasksUsecase(taskRepo)
	sessionRepo := repository.InitInMemorySessionRepository()
	sessionsUsecase := sessions.InitSessionsUsecase(sessionRepo)
	engine := gin.Default()
	routes.AddRoutes(engine, tasksUsecase, sessionsUsecase)
	engine.Run()
}
