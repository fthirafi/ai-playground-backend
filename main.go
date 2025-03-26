package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Initialize Fiber app
	app := fiber.New()

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if (err != nil) {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	// Database and collection
	db := client.Database("tradomlab")
	reservations := db.Collection("reservations")

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to TradomLab Meeting Room Reservation API")
	})

	app.Get("/reservations", func(c *fiber.Ctx) error {
		cursor, err := reservations.Find(context.TODO(), bson.M{})
		if err != nil {
			return c.Status(500).SendString("Failed to fetch reservations")
		}
		defer cursor.Close(context.TODO())

		var results []bson.M
		if err := cursor.All(context.TODO(), &results); err != nil {
			return c.Status(500).SendString("Failed to parse reservations")
		}

		return c.JSON(results) // Returns room, time, reservedBy, and company
	})

	// Start server
	log.Fatal(app.Listen(":3000"))
}
