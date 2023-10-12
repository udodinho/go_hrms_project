package employee

import (
	"github.com/gofiber/fiber/v2"
	"github.com/udodinho/hrms/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Employee struct {
	ID     string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Age    float64 `json:"age"`
}

func GetEmployees(c *fiber.Ctx) error {
	query := bson.D{{}}
	cursor, err := database.MG.Db.Collection("employees").Find(c.Context(), query)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	var employees []Employee = make([]Employee, 0)

	if err := cursor.All(c.Context(), &employees); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(employees)
	
}

func GetEmployee(c *fiber.Ctx) error {
	var employee Employee
	collection := database.MG.Db.Collection("employees")
	id := c.Params("id")

	employeeID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	query := bson.D{{Key: "_id", Value: employeeID}}

	err = collection.FindOne(c.Context(), query).Decode(&employee)

	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.SendStatus(400)
		}
		return c.SendStatus(500)
	}

	return c.Status(200).JSON(employee)

}

func CreateEmployee(c *fiber.Ctx) error {
	collection := database.MG.Db.Collection("employees")

	employee := new(Employee)

	if err := c.BodyParser(employee); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	employee.ID = ""

	insertionResult, err := collection.InsertOne(c.Context(), employee)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}

	createdRecord := collection.FindOne(c.Context(), filter)
	createdEmployee := &Employee{}
	createdRecord.Decode(createdEmployee)

	return c.Status(201).JSON(createdEmployee)

}

func UpdateEmployee(c *fiber.Ctx) error {
	collection := database.MG.Db.Collection("employees")
	id := c.Params("id")

	employeeID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	employee := new(Employee)

	if err := c.BodyParser(employee); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	query := bson.D{{Key: "_id", Value: employeeID}}

	updateEmployee := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "name", Value: employee.Name},
				{Key: "salary", Value: employee.Salary},
				{Key: "age", Value: employee.Age},
			},
		},
	}

	err = collection.FindOneAndUpdate(c.Context(), query, updateEmployee).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.SendStatus(400)
		}
		return c.SendStatus(500)
	}

	employee.ID = id

	return c.Status(200).JSON(employee)

}

func DeleteEmployee(c *fiber.Ctx) error {
	collection := database.MG.Db.Collection("employees")
	id := c.Params("id")

	employeeID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	query := bson.D{{Key: "_id", Value: employeeID}}

	result, err := collection.DeleteOne(c.Context(), &query)

	if err != nil {
		return c.SendStatus(500)
	}

	if result.DeletedCount < 1 {
		return c.Status(404).JSON("No employee with the id")
	}

	return c.Status(200).JSON("Record deleted successfully")
}
