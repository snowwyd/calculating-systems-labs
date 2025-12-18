package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NotificationRepository определяет интерфейс для работы с уведомлениями
type NotificationRepository interface {
	Create(ctx context.Context, notification *Notification) error
	FindAll(ctx context.Context) ([]*Notification, error)
}

// MongoDBNotificationRepository реализует репозиторий для MongoDB
type MongoDBNotificationRepository struct {
	collection *mongo.Collection
}

// NewMongoDBNotificationRepository создаёт новый репозиторий
func NewMongoDBNotificationRepository(collection *mongo.Collection) *MongoDBNotificationRepository {
	return &MongoDBNotificationRepository{collection: collection}
}

// Create создаёт новое уведомление
func (r *MongoDBNotificationRepository) Create(ctx context.Context, notification *Notification) error {
	_, err := r.collection.InsertOne(ctx, notification)
	if err != nil {
		return ErrDatabaseError{Operation: "создание уведомления", Err: err}
	}
	return nil
}

// FindAll возвращает все уведомления
func (r *MongoDBNotificationRepository) FindAll(ctx context.Context) ([]*Notification, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, ErrDatabaseError{Operation: "получение всех уведомлений", Err: err}
	}
	defer cursor.Close(ctx)

	notifications := []*Notification{}
	for cursor.Next(ctx) {
		var notification Notification
		if err := cursor.Decode(&notification); err != nil {
			return nil, ErrDatabaseError{Operation: "декодирование уведомления", Err: err}
		}
		notifications = append(notifications, &notification)
	}

	if err := cursor.Err(); err != nil {
		return nil, ErrDatabaseError{Operation: "итерация по уведомлениям", Err: err}
	}

	return notifications, nil
}

// InitDatabase инициализирует базу данных
func InitDatabase(uri string) (*mongo.Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	collection := client.Database("lab2db").Collection("notifications")
	return collection, nil
}

