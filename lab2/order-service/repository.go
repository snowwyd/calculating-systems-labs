package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// OrderRepository определяет интерфейс для работы с заказами
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	FindByCode(ctx context.Context, code string) (*Order, error)
	FindAll(ctx context.Context) ([]*Order, error)
	Delete(ctx context.Context, code string) error
}

// MongoDBOrderRepository реализует репозиторий для MongoDB
type MongoDBOrderRepository struct {
	collection *mongo.Collection
}

// NewMongoDBOrderRepository создаёт новый репозиторий
func NewMongoDBOrderRepository(collection *mongo.Collection) *MongoDBOrderRepository {
	return &MongoDBOrderRepository{collection: collection}
}

// Create создаёт новый заказ
func (r *MongoDBOrderRepository) Create(ctx context.Context, order *Order) error {
	_, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return ErrDatabaseError{Operation: "создание заказа", Err: err}
	}
	return nil
}

// FindByCode находит заказ по коду
func (r *MongoDBOrderRepository) FindByCode(ctx context.Context, code string) (*Order, error) {
	var order Order
	err := r.collection.FindOne(ctx, bson.M{"order_code": code}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrOrderNotFound{Code: code}
		}
		return nil, ErrDatabaseError{Operation: "поиск заказа", Err: err}
	}
	return &order, nil
}

// FindAll возвращает все заказы
func (r *MongoDBOrderRepository) FindAll(ctx context.Context) ([]*Order, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, ErrDatabaseError{Operation: "получение всех заказов", Err: err}
	}
	defer cursor.Close(ctx)

	orders := []*Order{}
	for cursor.Next(ctx) {
		var order Order
		if err := cursor.Decode(&order); err != nil {
			return nil, ErrDatabaseError{Operation: "декодирование заказа", Err: err}
		}
		orders = append(orders, &order)
	}

	if err := cursor.Err(); err != nil {
		return nil, ErrDatabaseError{Operation: "итерация по заказам", Err: err}
	}

	return orders, nil
}

// Delete удаляет заказ
func (r *MongoDBOrderRepository) Delete(ctx context.Context, code string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"order_code": code})
	if err != nil {
		return ErrDatabaseError{Operation: "удаление заказа", Err: err}
	}

	if result.DeletedCount == 0 {
		return ErrOrderNotFound{Code: code}
	}

	return nil
}

// InitDatabase инициализирует базу данных
func InitDatabase(uri string) (*mongo.Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	collection := client.Database("lab2db").Collection("orders")
	return collection, nil
}

