package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Notification struct {
	MessageType string `json:"message_type" bson:"message_type"`
	Description string `json:"description" bson:"description"`
	Date        string `json:"date" bson:"date"`
}

var collection *mongo.Collection

func initDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:password@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("lab2db").Collection("notifications")
}

func createNotification(w http.ResponseWriter, r *http.Request) {
	var notif Notification
	json.NewDecoder(r.Body).Decode(&notif)

	if notif.Date == "" {
		notif.Date = time.Now().Format("2006-01-02 15:04:05")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, notif)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notif)
}

func getNotifications(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	notifications := []Notification{}
	for cursor.Next(ctx) {
		var notif Notification
		cursor.Decode(&notif)
		notifications = append(notifications, notif)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func main() {
	initDB()

	r := mux.NewRouter()
	r.HandleFunc("/notifications", createNotification).Methods("POST")
	r.HandleFunc("/notifications", getNotifications).Methods("GET")

	log.Println("Notification Service запущен на :8084")
	log.Fatal(http.ListenAndServe(":8084", r))
}

