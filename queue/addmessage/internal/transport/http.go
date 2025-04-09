package transport

import (
	"log/slog"
	"main/internal/buisness"
	"main/internal/config"
	"net/http"

	"github.com/rs/cors"
)

type Hundl struct {
	Mess buisness.Message
}

func Hudlers(cfg *config.Config, sl *slog.Logger) {

	var h Hundl

	h.Mess.Config = cfg
	h.Mess.Logger = sl

	mux := http.NewServeMux()

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	// Подключаем проект chatapp
	mux.HandleFunc("/mango_internal/addmessage", h.Mess.Add)

	handler := cors.Handler(mux)
	http.ListenAndServe(":50551", handler)

}
