package postgres

import (
	"context"
	"errors"
	"time"

	"myapp/internal/models"

	"github.com/jmoiron/sqlx"
)

// Repository взаимодействует с базой данных для операций с группами
type Repository struct {
	db *sqlx.DB
}

// NewRepository создает новый экземпляр репозитория
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// GetGroupMembers получает всех участников зала
func (r *Repository) GetGroupMembers(ctx context.Context, gymID string) ([]models.User, error) {
	var users []models.User

	query := `
		SELECT u.id, u.email, u.first_name, u.last_name, gm.status, u.created_at, u.updated_at
		FROM users u
		JOIN group_members gm ON u.id = gm.user_id
		WHERE gm.gym_id = $1
	`

	err := r.db.SelectContext(ctx, &users, query, gymID)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserGroup получает группу, к которой принадлежит пользователь
func (r *Repository) GetUserGroup(ctx context.Context, userID string) (models.Group, []models.User, error) {
	var group models.Group
	var users []models.User

	// Сначала получаем зал, к которому принадлежит пользователь
	queryGym := `
		SELECT g.id, g.gym_id, g.name, g.created_at
		FROM groups g
		JOIN group_members gm ON g.gym_id = gm.gym_id
		WHERE gm.user_id = $1
		LIMIT 1
	`

	err := r.db.GetContext(ctx, &group, queryGym, userID)
	if err != nil {
		return models.Group{}, nil, err
	}

	// Затем получаем всех участников этого зала
	users, err = r.GetGroupMembers(ctx, group.GymID)
	if err != nil {
		return models.Group{}, nil, err
	}

	return group, users, nil
}

// UpdateUserStatus обновляет статус активности пользователя в зале
func (r *Repository) UpdateUserStatus(ctx context.Context, userID, gymID string, status models.ActivityStatus) error {
	query := `
		UPDATE group_members
		SET status = $1, updated_at = $2
		WHERE user_id = $3 AND gym_id = $4
		RETURNING id
	`

	var id string
	err := r.db.QueryRowContext(ctx, query, status, time.Now(), userID, gymID).Scan(&id)
	if err != nil {
		return err
	}

	if id == "" {
		return errors.New("пользователь не найден в этом зале")
	}

	return nil
}

// GetUserStatus получает статус пользователя в конкретном зале
func (r *Repository) GetUserStatus(ctx context.Context, userID, gymID string) (models.ActivityStatus, error) {
	query := `
		SELECT status
		FROM group_members
		WHERE user_id = $1 AND gym_id = $2
	`

	var status models.ActivityStatus
	err := r.db.GetContext(ctx, &status, query, userID, gymID)
	if err != nil {
		return "", err
	}

	return status, nil
}

// AddUserToGym добавляет пользователя в группу зала
func (r *Repository) AddUserToGym(ctx context.Context, userID, gymID string) error {
	query := `
		INSERT INTO group_members (user_id, gym_id, status, joined_at, updated_at)
		VALUES ($1, $2, $3, $4, $4)
		ON CONFLICT (user_id, gym_id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query, userID, gymID, models.ActiveStatus, time.Now())
	return err
}