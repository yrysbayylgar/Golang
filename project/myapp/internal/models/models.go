package models

import (
	"time"
)

// ActivityStatus представляет статус активности пользователя в зале
type ActivityStatus string

// Константы статуса активности
const (
	ActiveStatus   ActivityStatus = "active"   // активен - ходит в зал
	InactiveStatus ActivityStatus = "inactive" // неактивен - не ходит в зал
)

// User представляет пользователя в сервисе групп
type User struct {
	ID        string         `json:"id" db:"id"`
	Email     string         `json:"email" db:"email"`
	FirstName string         `json:"first_name" db:"first_name"`
	LastName  string         `json:"last_name" db:"last_name"`
	Status    ActivityStatus `json:"status" db:"status"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

// GroupMember представляет членство пользователя в группе зала
type GroupMember struct {
	ID        string         `json:"id" db:"id"`
	UserID    string         `json:"user_id" db:"user_id"`
	GymID     string         `json:"gym_id" db:"gym_id"`
	Status    ActivityStatus `json:"status" db:"status"`
	JoinedAt  time.Time      `json:"joined_at" db:"joined_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

// Group представляет группу пользователей в зале
type Group struct {
	ID        string    `json:"id" db:"id"`
	GymID     string    `json:"gym_id" db:"gym_id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// UpdateStatusRequest представляет запрос на обновление статуса пользователя
type UpdateStatusRequest struct {
	Status ActivityStatus `json:"status"`
}

// GroupMembersResponse представляет ответ со списком участников группы
type GroupMembersResponse struct {
	Members []User `json:"members"`
}
