package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Maden-in-haven/crmlib/pkg/database"
)

type UniversalRequest struct {
	ID          string
	Username    string `json:"username" validate:"required"` // Имя пользователя (обязательно)
	Role        string
	FullName    *string                 `json:"full_name,omitempty"`                               // Полное имя пользователя
	PhoneNumber *string                 `json:"phone_number,omitempty" validate:"omitempty,phone"` // Проверяем, если указано, что номер телефона валиден и начинается с +7
	HireDate    *string                 `json:"hire_date,omitempty" validate:"omitempty,rfc3339"`  // Проверяем, если указано, что дата в формате RFC3339          // Дата найма (только для менеджеров)
	Permissions *map[string]interface{} `json:"permissions,omitempty"`                             // Права доступа (только для администраторов)
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	user, err := database.DB.GetUserByID(context.Background(), userID)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	universalRequest := UniversalRequest{
		ID:       userID,
		Username: user.Username,
		Role:     user.Role,
	}
	if user.Role == "Admin" {
		admin, err := database.DB.GetAdminByID(context.Background(), userID)
		if err != nil {
			http.Error(w, "Ошибка получения данных администратора", http.StatusInternalServerError)
			return
		}
		universalRequest.Permissions = &admin.Permissions

	} else if user.Role == "Manager" {
		manager, err := database.DB.GetManagerByID(context.Background(), userID)
		if err != nil {
			http.Error(w, "Ошибка получения данных менеджера", http.StatusInternalServerError)
			return
		}
		universalRequest.FullName = &manager.FullName
		universalRequest.HireDate = &manager.HireDate
	} else if user.Role == "Client" {
		client, err := database.DB.GetClientByID(context.Background(), userID)
		if err != nil {
			http.Error(w, "Ошибка получения данных клиента", http.StatusInternalServerError)
			return
		}
		universalRequest.FullName = &client.FullName
		universalRequest.PhoneNumber = &client.PhoneNumber
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(universalRequest); err != nil {
		http.Error(w, "Ошибка кодирования ответа", http.StatusInternalServerError)
		return
	}
}
