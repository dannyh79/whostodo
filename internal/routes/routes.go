package routes

import (
	"net/http"

	"github.com/dannyh79/whostodo/internal/tasks"
	"github.com/gin-gonic/gin"
)

func AddRoutes(r *gin.Engine) *gin.Engine {
	r.GET("/tasks", func (c *gin.Context) {
		tasks := tasks.ListTasks()
		c.JSON(http.StatusOK, gin.H{
			"result": tasks,
		})
	})
	return r
}
