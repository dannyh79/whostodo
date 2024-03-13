package main

import (
	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/routes"
	"github.com/dannyh79/whostodo/internal/sessions"
	"github.com/dannyh79/whostodo/internal/tasks"
	"github.com/gin-gonic/gin"
)

func main() {
	repo := repository.InitInMemoryTaskRepository()
	tasksUsecase := tasks.InitTasksUsecase(repo)
	sessionsUsecase := sessions.InitSessionsUsecase()
	engine := gin.Default()
	routes.AddRoutes(engine, tasksUsecase, sessionsUsecase)
	engine.Run()
}
