package main

// Notification представляет уведомление
type Notification struct {
	MessageType string `json:"message_type" bson:"message_type"`
	Description string `json:"description" bson:"description"`
	Date        string `json:"date" bson:"date"`
}

