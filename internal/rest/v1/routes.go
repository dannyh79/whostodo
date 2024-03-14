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
type ListTaskItem struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}
type ListTasksOutput struct {
	Result []ListTaskItem `json:"result"`
}

type PostTaskOutput struct {
	Result struct {
		Name   string `json:"name"`
		Status int    `json:"status"`
		Id     int    `json:"id"`
	} `json:"result"`
}

type UpdateTaskOutput struct {
	Result struct {
		Name   string `json:"name"`
		Status int    `json:"status"`
		Id     int    `json:"id"`
	} `json:"result"`
}

type FailedUpdateTaskOutput struct {
	Result struct{} `json:"result"`
}

type PostAuthSuccessOutput struct {
	Token string `json:"result"`
}

type PostAuthNotModifiedOutput struct{}

var UnprotectedPaths = map[string]string{
	"auth": "/auth",
}

func AddRoutes(r *gin.Engine, tasksU *tasks.TasksUsecase, sessionsU *sessions.SessionsUsecase) {
	v1 := r.Group("/v1")

	v1.Use(sessionMiddleware(sessionsU, UnprotectedPaths))

	v1.POST(UnprotectedPaths["auth"], authenticateHandler(sessionsU))

	v1.GET("/tasks", listTasksHandler(tasksU))
	v1.POST("/task", createTaskHandler(tasksU))
	v1.PUT("/task/:id", updateTaskHandler(tasksU))
	v1.DELETE("/task/:id", deleteTaskHandler(tasksU))
}

func listTasksHandler(u *tasks.TasksUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		tasks := u.ListTasks()
		c.JSON(http.StatusOK, toListTasksOutput(tasks))
	}
}

func createTaskHandler(u *tasks.TasksUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload tasks.CreateTaskInput
		c.ShouldBind(&payload)
		task := u.CreateTask(&payload)
		c.JSON(http.StatusCreated, toPostTaskOutput(task))
	}
}

func updateTaskHandler(u *tasks.TasksUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var payload tasks.UpdateTaskInput
		c.ShouldBind(&payload)

		updated, err := u.UpdateTask(id, &payload)
		if err != nil {
			c.JSON(http.StatusNotFound, FailedUpdateTaskOutput{})
			return
		}

		c.JSON(http.StatusCreated, toUpdateTaskOutput(updated))
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
		token := getTokenFromHeader(c)
		if u.Validate(token) {
			c.JSON(http.StatusNotModified, PostAuthNotModifiedOutput{})
			return
		}

		token = u.Authenticate()
		c.JSON(http.StatusCreated, PostAuthSuccessOutput{Token: token})
	}
}

func sessionMiddleware(u *sessions.SessionsUsecase, ignore map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, path := range ignore {
			if c.Request.URL.Path == "/v1"+path {
				c.Next()
				return
			}
		}

		token := getTokenFromHeader(c)
		if !u.Validate(token) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{})
			return
		}

		c.Next()
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

func toListTasksOutput(ts []*tasks.TaskOutput) *ListTasksOutput {
	var result = make([]ListTaskItem, 0)
	var output ListTasksOutput
	for _, t := range ts {
		result = append(result, ListTaskItem{
			Id:     t.Id,
			Name:   t.Name,
			Status: t.Status,
		})
	}
	output.Result = result
	return &output
}

func toPostTaskOutput(t *tasks.TaskOutput) *PostTaskOutput {
	var output PostTaskOutput
	output.Result.Id = t.Id
	output.Result.Name = t.Name
	output.Result.Status = t.Status
	return &output
}

func toUpdateTaskOutput(t *tasks.TaskOutput) *UpdateTaskOutput {
	var output UpdateTaskOutput
	output.Result.Id = t.Id
	output.Result.Name = t.Name
	output.Result.Status = t.Status
	return &output
}
