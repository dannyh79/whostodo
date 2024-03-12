package main

import (
	routes "github.com/dannyh79/whostodo/internal"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	routes.AddRoutes(engine)
	engine.Run()
}
