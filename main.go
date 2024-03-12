package main

import (
	"github.com/dannyh79/whostodo/internal/repository"
	"github.com/dannyh79/whostodo/internal/routes"
	"github.com/dannyh79/whostodo/internal/tasks"
	"github.com/gin-gonic/gin"
)

func main() {
	repo := repository.InitInMemoryTaskRepository()
	usecase := tasks.InitTasksUsecase(repo)
	engine := gin.Default()
	routes.AddRoutes(engine, usecase)
	engine.Run()
}
