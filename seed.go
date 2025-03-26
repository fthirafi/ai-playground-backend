package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func generateTimeslots() []string {
	timeslots := []string{}
	for hour := 9; hour < 17; hour++ {
		timeslots = append(timeslots, 
			formatTimeslot(hour, 0, hour, 30),
			formatTimeslot(hour, 30, hour+1, 0),
		)
	}
	return timeslots
}

func formatTimeslot(startHour, startMinute, endHour, endMinute int) string {
	return formatTime(startHour, startMinute) + " - " + formatTime(endHour, endMinute)
}

func formatTime(hour, minute int) string {
	period := "AM"
	if hour >= 12 {
		period = "PM"
		if hour > 12 {
			hour -= 12
		}
	}
	return time.Date(0, 0, 0, hour, minute, 0, 0, time.UTC).Format("3:04") + " " + period
}

func main() {
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	// Database and collection
	db := client.Database("tradomlab")
	reservations := db.Collection("reservations")

	// Generate timeslots
	timeslots := generateTimeslots()

	// Seed data
	initialData := []interface{}{
		bson.M{"room": "Room A", "time": timeslots[0], "reservedBy": "John Doe", "company": "TechCorp"},
		bson.M{"room": "Room B", "time": timeslots[3], "reservedBy": "Jane Smith", "company": "Innovate Inc"},
		bson.M{"room": "Room C", "time": timeslots[5], "reservedBy": "Alice Brown", "company": "Alpha Solutions"},
		bson.M{"room": "Room A", "time": timeslots[7], "reservedBy": "Bob Johnson", "company": "Beta Tech"},
		bson.M{"room": "Room B", "time": timeslots[9], "reservedBy": "Emily Davis", "company": "Gamma Group"},
	}

	_, err = reservations.InsertMany(context.TODO(), initialData)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database seeded successfully!")
}
