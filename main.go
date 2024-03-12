package main

import (
	"github.com/dannyh79/whostodo/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	routes.AddRoutes(engine)
	engine.Run()
}
