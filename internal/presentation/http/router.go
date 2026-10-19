package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"mg-back/internal/presentation/http/handlers"
	"mg-back/internal/presentation/http/middleware"
	"mg-back/internal/presentation/http/response"
)

func NewRouter(
	authHandler *handlers.AuthHandler,
	passwordResetHandler *handlers.PasswordResetHandler,
	scaleHandler *handlers.ScaleHandler,
	gameHandler *handlers.GameHandler,
	leaderboardHandler *handlers.LeaderboardHandler,
	jwtSecret []byte,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
    AllowedOrigins: []string{
        "https://musicalgame.netlify.app",
        "http://localhost:5173",
    },
    AllowedMethods: []string{
        "GET",
        "POST",
        "OPTIONS",
    },
    AllowedHeaders: []string{
        "Accept",
        "Authorization",
        "Content-Type",
    },
    AllowCredentials: false,
    MaxAge: 300,
	}))

	r.Get("/health", healthCheck)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.Refresh)
		r.Post("/forgot-password", passwordResetHandler.Forgot)
		r.Post("/verify-reset-code", passwordResetHandler.VerifyCode)
		r.Post("/reset-password", passwordResetHandler.Reset)
	})

	// Catálogo de juego: lectura pública
	r.Get("/scales", scaleHandler.List)
	r.Get("/scales/{scaleID}/levels", scaleHandler.ListLevels)
	r.Get("/levels/{levelID}", scaleHandler.GetLevel)
	r.Get("/leaderboard/{levelID}", leaderboardHandler.Get)

	// Rutas protegidas
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(jwtSecret))
		r.Post("/attempts", gameHandler.SubmitAttempt)
	})

	return r
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
