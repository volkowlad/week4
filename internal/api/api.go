package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"week4/internal/api/mw"
	"week4/internal/service"
)

type Routers struct {
	Service service.Service
}

// NewRouters - конструктор для настройки API
func NewRouters(r *Routers, token string) *fiber.App {
	app := fiber.New()

	// Настройка CORS (разрешенные методы, заголовки, авторизация)
	app.Use(cors.New(cors.Config{
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowHeaders:     "Accept, Authorization, Content-Type, X-CSRF-Token, X-REQUEST-SomeID",
		ExposeHeaders:    "Link",
		AllowCredentials: false, // change
		MaxAge:           300,
	}))

	// Группа маршрутов с авторизацией
	apiGroup := app.Group("/v1", mw.Authorization(token))

	// Роут для создания задачи
	{
		apiGroup.Post("/create_task", r.Service.CreateTask)
		apiGroup.Get("/task/:id", r.Service.GetTask)
		apiGroup.Get("/tasks", r.Service.GetAllTasks)
		apiGroup.Delete("/delete/:id", r.Service.DeleteTask)
		apiGroup.Put("/update/:id", r.Service.UpdateTask)
	}

	return app
}
