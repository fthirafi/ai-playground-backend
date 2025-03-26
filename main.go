package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Reservation struct {
	Room       string `json:"room"`
	Time       string `json:"time"`
	ReservedBy string `json:"reservedBy"`
	Company    string `json:"company"`
}

var reservationsCollection *mongo.Collection

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getReservations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cursor, err := reservationsCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch reservations", http.StatusInternalServerError)
		log.Printf("Error fetching reservations: %v", err)
		return
	}
	defer cursor.Close(context.TODO())

	var results []Reservation
	if err := cursor.All(context.TODO(), &results); err != nil {
		http.Error(w, "Failed to parse reservations", http.StatusInternalServerError)
		log.Printf("Error parsing reservations: %v", err)
		return
	}

	if len(results) == 0 {
		http.Error(w, "No reservations found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Failed to encode reservations", http.StatusInternalServerError)
		log.Printf("Error encoding reservations: %v", err)
	}
}

func createReservation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newReservation Reservation
	if err := json.NewDecoder(r.Body).Decode(&newReservation); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	_, err := reservationsCollection.InsertOne(context.TODO(), newReservation)
	if err != nil {
		http.Error(w, "Failed to create reservation", http.StatusInternalServerError)
		log.Printf("Error inserting reservation: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newReservation); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

func main() {
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	// Initialize collection
	db := client.Database("tradomlab")
	reservationsCollection = db.Collection("reservations")

	mux := http.NewServeMux()
	mux.HandleFunc("/reservations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getReservations(w, r)
		} else if r.Method == http.MethodPost {
			createReservation(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Starting server on port 3000...")
	if err := http.ListenAndServe(":3000", enableCORS(mux)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
