package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Task struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Status int `json:"status"`
}

func Server() *gin.Engine {
	engine := gin.Default()
	engine.GET("/tasks", func (c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"result": []Task{{Id: 1, Name: "name", Status: 0}},
		})
	})
	return engine
}
