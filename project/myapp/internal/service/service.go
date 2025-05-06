package service

import (
	"context"
	"errors"

	"myapp/internal/models"
)

// Repository определяет интерфейс для операций с базой данных
type Repository interface {
	GetGroupMembers(ctx context.Context, gymID string) ([]models.User, error)
	GetUserGroup(ctx context.Context, userID string) (models.Group, []models.User, error)
	UpdateUserStatus(ctx context.Context, userID, gymID string, status models.ActivityStatus) error
	GetUserStatus(ctx context.Context, userID, gymID string) (models.ActivityStatus, error)
	AddUserToGym(ctx context.Context, userID, gymID string) error
}

// Service обрабатывает бизнес-логику для сервиса групп
type Service struct {
	repo Repository
}

// NewService создает новый сервис групп
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetGroupMembers получает всех участников зала
func (s *Service) GetGroupMembers(ctx context.Context, gymID string) ([]models.User, error) {
	if gymID == "" {
		return nil, errors.New("требуется ID зала")
	}

	return s.repo.GetGroupMembers(ctx, gymID)
}

// GetUserGroup получает информацию о группе и участниках для пользователя
func (s *Service) GetUserGroup(ctx context.Context, userID string) (models.Group, []models.User, error) {
	if userID == "" {
		return models.Group{}, nil, errors.New("требуется ID пользователя")
	}

	return s.repo.GetUserGroup(ctx, userID)
}

// UpdateUserStatus обновляет статус пользователя в зале
func (s *Service) UpdateUserStatus(ctx context.Context, userID, gymID string, status models.ActivityStatus) error {
	if userID == "" || gymID == "" {
		return errors.New("требуются ID пользователя и ID зала")
	}

	// Проверка статуса
	if status != models.ActiveStatus && status != models.InactiveStatus {
		return errors.New("недопустимое значение статуса")
	}

	return s.repo.UpdateUserStatus(ctx, userID, gymID, status)
}

// GetUserStatus получает статус пользователя в зале
func (s *Service) GetUserStatus(ctx context.Context, userID, gymID string) (models.ActivityStatus, error) {
	if userID == "" || gymID == "" {
		return "", errors.New("требуются ID пользователя и ID зала")
	}

	return s.repo.GetUserStatus(ctx, userID, gymID)
}

// AddUserToGym добавляет пользователя в группу зала
func (s *Service) AddUserToGym(ctx context.Context, userID, gymID string) error {
	if userID == "" || gymID == "" {
		return errors.New("требуются ID пользователя и ID зала")
	}

	return s.repo.AddUserToGym(ctx, userID, gymID)
}
