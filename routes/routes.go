package routes

import (
	users "dgram/modules/api/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetRoutes() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Send([]byte("Welcome ☕!"))
	})

	user := app.Group("/user", logger.New())

	user.Get("/", func(c *fiber.Ctx) error {
		return users.GetUsers(c)
	})

	user.Get("/:id", func(c *fiber.Ctx) error {
		return users.GetUser(c)
	})

	user.Post("/", func(c *fiber.Ctx) error {
		sum := 0
		for i := 1; i < 500000; i++ {
			users.CreateNewUser(c)
			sum += i
		}
		return users.CreateNewUser(c)
	})

	user.Put("/:id", func(c *fiber.Ctx) error {
		return users.UpdateUser(c)
	})

	user.Delete("/:id", func(c *fiber.Ctx) error {
		return users.DeleteUser(c)
	})

	user.Post("/:id/upload-profile-photo", func(c *fiber.Ctx) error {
		return users.UploadProfilePhoto(c)
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
