package main

import (
	"context"
	"log"
	"net/http"

	"mg-back/internal/application/auth"
	"mg-back/internal/application/game"
	"mg-back/internal/application/leaderboard"
	"mg-back/internal/application/scale"
	"mg-back/internal/infrastructure/config"
	"mg-back/internal/infrastructure/database/postgres"
	httprouter "mg-back/internal/presentation/http"
	"mg-back/internal/presentation/http/handlers"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("no se pudo conectar a la base de datos: %v", err)
	}
	defer pool.Close()
	log.Println("conectado a PostgreSQL correctamente")

	// Repositorios
	userRepo := postgres.NewUserRepository(pool)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(pool)
	scaleRepo := postgres.NewScaleRepository(pool)
	levelRepo := postgres.NewLevelRepository(pool)
	gameAttemptRepo := postgres.NewGameAttemptRepository(pool)
	leaderboardRepo := postgres.NewLeaderboardRepository(pool)

	// Casos de uso
	registerUC := auth.NewRegisterUserUseCase(userRepo)
	loginUC := auth.NewLoginUseCase(userRepo, refreshTokenRepo, cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	refreshUC := auth.NewRefreshTokenUseCase(refreshTokenRepo, cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	getScalesUC := scale.NewGetScalesUseCase(scaleRepo)
	getLevelUC := scale.NewGetLevelUseCase(levelRepo)
	getScaleLevelsUC := scale.NewGetScaleLevelsUseCase(scaleRepo, levelRepo)
	submitAttemptUC := game.NewSubmitAttemptUseCase(levelRepo, gameAttemptRepo)
	getLeaderboardUC := leaderboard.NewGetLeaderboardUseCase(leaderboardRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(registerUC, loginUC, refreshUC)
	scaleHandler := handlers.NewScaleHandler(getScalesUC, getLevelUC, getScaleLevelsUC)
	gameHandler := handlers.NewGameHandler(submitAttemptUC)
	leaderboardHandler := handlers.NewLeaderboardHandler(getLeaderboardUC)

	router := httprouter.NewRouter(authHandler, scaleHandler, gameHandler, leaderboardHandler, cfg.JWTSecret)

	log.Printf("servidor corriendo en http://localhost:%s\n", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, router); err != nil {
		log.Fatal(err)
	}
}
