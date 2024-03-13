package routes

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/dannyh79/whostodo/internal/sessions"
	"github.com/dannyh79/whostodo/internal/tasks"
	"github.com/gin-gonic/gin"
)

// Returning format is slightly different per spec
type PostTaskOutput struct {
	Name   string `json:"name"`
	Status int    `json:"status"`
	Id     int    `json:"id"`
}

var UnprotectedPaths = map[string]string{
	"auth": "/auth",
}

func AddRoutes(r *gin.Engine, tasksU *tasks.TasksUsecase, sessionsU *sessions.SessionsUsecase) {
	r.Use(sessionMiddleware(sessionsU, UnprotectedPaths))

	r.POST(UnprotectedPaths["auth"], authenticateHandler(sessionsU))

	r.GET("/tasks", listTasksHandler(tasksU))
	r.POST("/task", createTaskHandler(tasksU))
	r.PUT("/task/:id", updateTaskHandler(tasksU))
	r.DELETE("/task/:id", deleteTaskHandler(tasksU))
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

		updated, err := u.UpdateTask(id, &payload)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"result": gin.H{},
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"result": toPostTaskOutput(updated),
		})
	}
}

func deleteTaskHandler(u *tasks.TasksUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		err := u.DeleteTask(id)
		if err != nil {
			c.JSON(http.StatusNotFound, nil)
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func authenticateHandler(u *sessions.SessionsUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	}
}

func sessionMiddleware(u *sessions.SessionsUsecase, ignore map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, path := range ignore {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		token := getTokenFromHeader(c)
		if u.Validate(token) {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{})
			return
		}
	}
}

func getTokenFromHeader(c *gin.Context) string {
	headerValue := c.Request.Header.Get("Authorization")
	bearerAndToken := strings.Split(headerValue, "Bearer ")
	if len(bearerAndToken) < 2 {
		return ""
	}

	return bearerAndToken[1]
}

func toPostTaskOutput(t *tasks.TaskOutput) PostTaskOutput {
	return PostTaskOutput{
		Id:     t.Id,
		Name:   t.Name,
		Status: t.Status,
	}
}
