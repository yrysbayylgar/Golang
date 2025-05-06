package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"myapp/internal/models"
	"myapp/pkg/auth"
	httputil "myapp/pkg/http"
)

// Service определяет интерфейс для бизнес-логики
type Service interface {
	GetGroupMembers(ctx context.Context, gymID string) ([]models.User, error)
	GetUserGroup(ctx context.Context, userID string) (models.Group, []models.User, error)
	UpdateUserStatus(ctx context.Context, userID, gymID string, status models.ActivityStatus) error
	GetUserStatus(ctx context.Context, userID, gymID string) (models.ActivityStatus, error)
	AddUserToGym(ctx context.Context, userID, gymID string) error
}

// Handler обрабатывает HTTP-запросы
type Handler struct {
	service Service
}

// NewHandler создает новый обработчик
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes регистрирует маршруты обработчика
func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/groups/{gymId}/members", h.GetGroupMembers).Methods("GET")
	r.HandleFunc("/groups/my", h.GetMyGroup).Methods("GET")
	r.HandleFunc("/groups/{gymId}/members/{userId}/status", h.GetUserStatus).Methods("GET")
	r.HandleFunc("/groups/{gymId}/members/{userId}/status", h.UpdateUserStatus).Methods("PUT")
	r.HandleFunc("/groups/{gymId}/members", h.AddUserToGym).Methods("POST")
}

// GetGroupMembers обрабатывает получение участников зала
func (h *Handler) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gymID := vars["gymId"]

	if gymID == "" {
		httputil.RespondWithError(w, http.StatusBadRequest, "Требуется ID зала")
		return
	}

	users, err := h.service.GetGroupMembers(r.Context(), gymID)
	if err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, "Ошибка получения участников группы")
		return
	}

	response := models.GroupMembersResponse{
		Members: users,
	}

	httputil.RespondWithJSON(w, http.StatusOK, response)
}

// GetMyGroup обрабатывает получение своей группы пользователем
func (h *Handler) GetMyGroup(w http.ResponseWriter, r *http.Request) {
	// Получаем ID пользователя из JWT токена
	userID, err := auth.GetUserIDFromToken(r)
	if err != nil {
		httputil.RespondWithError(w, http.StatusUnauthorized, "Недействительный токен")
		return
	}

	group, members, err := h.service.GetUserGroup(r.Context(), userID)
	if err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, "Ошибка получения группы")
		return
	}

	response := struct {
		Group   models.Group `json:"group"`
		Members []models.User `json:"members"`
	}{
		Group:   group,
		Members: members,
	}

	httputil.RespondWithJSON(w, http.StatusOK, response)
}

// GetUserStatus обрабатывает получение статуса пользователя в зале
func (h *Handler) GetUserStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gymID := vars["gymId"]
	userID := vars["userId"]

	if gymID == "" || userID == "" {
		httputil.RespondWithError(w, http.StatusBadRequest, "Требуются ID зала и ID пользователя")
		return
	}

	status, err := h.service.GetUserStatus(r.Context(), userID, gymID)
	if err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, "Ошибка получения статуса пользователя")
		return
	}

	response := struct {
		Status models.ActivityStatus `json:"status"`
	}{
		Status: status,
	}

	httputil.RespondWithJSON(w, http.StatusOK, response)
}

// UpdateUserStatus обрабатывает обновление статуса пользователя в зале
func (h *Handler) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gymID := vars["gymId"]
	userID := vars["userId"]

	if gymID == "" || userID == "" {
		httputil.RespondWithError(w, http.StatusBadRequest, "Требуются ID зала и ID пользователя")
		return
	}

	var req models.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.RespondWithError(w, http.StatusBadRequest, "Недопустимое тело запроса")
		return
	}

	if err := h.service.UpdateUserStatus(r.Context(), userID, gymID, req.Status); err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, "Ошибка обновления статуса пользователя")
		return
	}

	httputil.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Статус успешно обновлен"})
}

// AddUserToGym обрабатывает добавление пользователя в зал
func (h *Handler) AddUserToGym(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gymID := vars["gymId"]

	if gymID == "" {
		httputil.RespondWithError(w, http.StatusBadRequest, "Требуется ID зала")
		return
	}

	// Получаем ID пользователя из JWT токена
	userID, err := auth.GetUserIDFromToken(r)
	if err != nil {
		httputil.RespondWithError(w, http.StatusUnauthorized, "Недействительный токен")
		return
	}

	if err := h.service.AddUserToGym(r.Context(), userID, gymID); err != nil {
		httputil.RespondWithError(w, http.StatusInternalServerError, "Ошибка добавления пользователя в зал")
		return
	}

	httputil.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Пользователь успешно добавлен в зал"})
}