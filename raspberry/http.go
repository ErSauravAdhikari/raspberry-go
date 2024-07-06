package rasberry

import (
	_ "github.com/ersauravadhikari/raspberry-go/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Raspberry API
// @version 1.0
// @description This is a simple task scheduler API.
// @BasePath /api/
// @securityDefinitions.apiKey ApiKeyAuth
// @in query
// @name api_key

// RunAPI starts the API server
// @Summary Start API server
// @Description Start the API server to manage tasks and schedules
// @Produce json
// @Success 200 {object} string "API server started"
// @Router / [get]
func (r *Raspberry) RunAPI(port string) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Load templates
	templates, err := loadTemplates()
	if err != nil {
		e.Logger.Fatal("Failed to load templates:", err)
	}
	e.Renderer = &Template{templates: templates}

	// Register routes for the login page and handler
	e.GET("/login", r.serveLoginPage)
	e.POST("/login", r.handleLogin)

	// Register routes for the web UI with cookie-based authentication
	web := e.Group("")

	web.Use(r.WebAuthMiddleware)

	web.GET("/", r.listTasks)
	web.GET("/task/:name", r.showTask)
	web.GET("/execution/:id", r.showExecution)

	// Register routes for the API with appropriate authentication
	api := e.Group("/api")
	api.Use(r.APIKeyAuthMiddleware)
	api.GET("/tasks", r.getTasks)
	api.GET("/task/:name/executions", r.getTaskExecutions)
	api.GET("/task_run/:id/logs", r.getTaskRunLogs)
	api.POST("/execution/:id/cancel", r.cancelExecutionByID)

	// Swagger docs endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":" + port))
}
