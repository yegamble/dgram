package routes

import (
	users "dgram/modules/api/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetRoutes() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("Helllo ☕!"))
	})

	user := app.Group("/user", logger.New())
	user.Get("/", func(c *fiber.Ctx) error {
		return users.CreateNewUser(c)
	})

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
