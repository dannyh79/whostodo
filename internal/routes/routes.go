package routes

import (
	"net/http"
	"strconv"

	"github.com/dannyh79/whostodo/internal/tasks"
	"github.com/gin-gonic/gin"
)

// Returning format is slightly different per spec
type PostTaskOutput struct {
	Name   string `json:"name"`
	Status int    `json:"status"`
	Id     int    `json:"id"`
}

func AddRoutes(r *gin.Engine, u *tasks.TasksUsecase) {
	r.GET("/tasks", listTasksHandler(u))
	r.POST("/tasks", createTaskHandler(u))
	r.PUT("/tasks/:id", updateTaskHandler(u))
}

func listTasksHandler(u *tasks.TasksUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		tasks := u.ListTasks()
		c.JSON(http.StatusOK, gin.H{
			"result": tasks,
		})
	}
}

func createTaskHandler(u *tasks.TasksUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload tasks.CreateTaskInput
		c.ShouldBind(&payload)
		task := u.CreateTask(&payload)
		c.JSON(http.StatusCreated, gin.H{
			"result": toPostTaskOutput(task),
		})
	}
}

func updateTaskHandler(u *tasks.TasksUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var payload tasks.UpdateTaskInput
		c.ShouldBind(&payload)

		updated := u.UpdateTask(id, &payload)
		c.JSON(http.StatusCreated, gin.H{
			"result": toPostTaskOutput(updated),
		})
	}
}

func toPostTaskOutput(t *tasks.TaskOutput) PostTaskOutput {
	return PostTaskOutput{
		Id:     t.Id,
		Name:   t.Name,
		Status: t.Status,
	}
}
