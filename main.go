package main

import (
	"dgram/database"
	users "dgram/modules/api/user"
	"dgram/modules/api/wallet"
	"dgram/routes"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

type Todo struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open(sqlite.Open("dgram.db"))
	if err != nil {
		panic("Failed to Connect to Database")
	}
	fmt.Println("Database connection successfully opened")

	database.DBConn.AutoMigrate(&users.User{})
	database.DBConn.AutoMigrate(&users.Post{})
	database.DBConn.AutoMigrate(&users.Comment{})
}

func main() {
	initDatabase()
	routes.SetUserRoutes()

	err := wallet.CheckTransactions()
	if err != nil {
		log.Fatal(err)
	}
}
