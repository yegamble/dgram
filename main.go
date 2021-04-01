package main

import (
	"dgram/database"
	users "dgram/modules/api/user"
	"dgram/routes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strconv"
)

type Todo struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

var todos = []*Todo{
	{Id: 1, Name: "Walk the Dog", Completed: false},
	{Id: 2, Name: "Walk the Cat", Completed: false},
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open(sqlite.Open("dgram.db"))
	if err != nil {
		panic("Failed to Connect to Database")
	}
	fmt.Println("Database connection successfully opened")

	database.DBConn.AutoMigrate(&users.User{})
}

func main() {
	initDatabase()
	routes.SetRoutes()
}

func UpdateTodo(ctx *fiber.Ctx) error {
	type request struct {
		Name      *string `json:"name"`
		Completed *bool   `json:"completed"`
	}

	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})

	}

	var body request
	err = ctx.BodyParser(&body)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse body",
		})
	}

	var todo *Todo

	for _, t := range todos {
		if t.Id == id {
			todo = t
			break
		}
	}

	if todo == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Not Found",
		})
	}

	if body.Name != nil {
		todo.Name = *body.Name
	}

	if body.Completed != nil {
		todo.Completed = *body.Completed
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}

func DeleteGetTodo(ctx *fiber.Ctx) error {
	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
	}

	for i, todo := range todos {
		if todo.Id == id {
			todos = append(todos[0:i], todos[i+1:]...)
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"success": true,
			})
		}
	}

	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "cannot find",
	})
}

func GetTodo(ctx *fiber.Ctx) error {
	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
	}

	for _, todo := range todos {
		if todo.Id == id {
			return ctx.Status(fiber.StatusOK).JSON(todo)
		}
	}

	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "cannot find todo",
	})
}

func CreateTodo(ctx *fiber.Ctx) error {
	type request struct {
		Name string `json:"name"`
	}

	var body request

	err := ctx.BodyParser(&body)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"error": "cannot parse request",
		})
	}

	todo := &Todo{
		Id:        len(todos) + 1,
		Name:      body.Name,
		Completed: false,
	}

	todos = append(todos, todo)
	return ctx.Status(fiber.StatusCreated).JSON(todo)
}

func GetTodos(ctx *fiber.Ctx) error {
	ctx.Status(fiber.StatusOK).JSON(todos)
	return nil
}
