package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetRoutes() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("Helllo ☕!"))
	})

	//api := app.Group("/todos", logger.New())

	//api.Get("/", GetTodos)
	//api.Post("/", CreateTodo)
	//api.Get("/:id", GetTodo)
	//api.Delete("/:id", DeleteGetTodo)
	//api.Patch("/:id", UpdateTodo)


	err := app.Listen("localhost:3000")
	if err != nil {
		panic(err)
	}
}
