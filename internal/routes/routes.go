package routes

import (
	"log"
	"net/http"
	"profile/internal/handler"
	"profile/internal/middlewares"

	"github.com/gorilla/mux"
)

func ProfileRun() {
	// Инициализируем новый маршрутизатор
	r := mux.NewRouter()

	r.Handle("api/profile", middlewares.JWTAuthentication(http.HandlerFunc(handler.GetProfile))).Methods("GET")

	// Запускаем сервер на порту 8080
	port := ":3000"
	log.Printf("Сервер запущен на порту %s", port)

	// Логируем ошибки, если таковые возникнут при запуске
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
