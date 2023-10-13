package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/udodinho/hrms/database"
	"github.com/udodinho/hrms/employee"
)

func init() {
	database.Connect()
}

func setUpRoutes(app *fiber.App) {
	api := app.Group("/api/v1/employee")
	api.Get("/", employee.GetEmployees)
	api.Post("/", employee.CreateEmployee)
	api.Get("/:id", employee.GetEmployee)
	api.Put("/:id", employee.UpdateEmployee)
	api.Delete("/:id", employee.DeleteEmployee)
}

func main() {
	app := fiber.New()
	port := ":3000"
	setUpRoutes(app)
	
	fmt.Printf("Server started on port %s", port)
	log.Fatal(app.Listen(port))
}
