package main

import (
	"context"
	"time"
)

// NotificationService предоставляет бизнес-логику для работы с уведомлениями
type NotificationService struct {
	repo NotificationRepository
}

// NewNotificationService создаёт новый сервис уведомлений
func NewNotificationService(repo NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

// CreateNotification создаёт новое уведомление
func (s *NotificationService) CreateNotification(ctx context.Context, req CreateNotificationRequest) (*Notification, error) {
	notification := &Notification{
		MessageType: req.MessageType,
		Description: req.Description,
		Date:        req.Date,
	}

	// Если дата не указана, устанавливаем текущую
	if notification.Date == "" {
		notification.Date = time.Now().Format("2006-01-02 15:04:05")
	}

	if err := s.repo.Create(ctx, notification); err != nil {
		return nil, err
	}

	return notification, nil
}

// GetAllNotifications получает все уведомления
func (s *NotificationService) GetAllNotifications(ctx context.Context) ([]*Notification, error) {
	return s.repo.FindAll(ctx)
}

