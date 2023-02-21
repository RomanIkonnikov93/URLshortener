package server

import (
	"net/http"

	"github.com/RomanIkonnikov93/URLshortner/cmd/config"
	"github.com/RomanIkonnikov93/URLshortner/internal/handlers"
	"github.com/RomanIkonnikov93/URLshortner/internal/repository"
	"github.com/RomanIkonnikov93/URLshortner/logging"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func StartServer(rep repository.Pool, cfg config.Config, logger logging.Logger) error {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(handlers.GzipRequest)
	r.Use(handlers.GzipResponse)
	r.Use(handlers.UserValidation(rep, logger))

	r.Post("/", handlers.PostHandler(rep, logger))
	r.Post("/api/shorten", handlers.PostJSONHandler(rep, logger))
	r.Post("/api/shorten/batch", handlers.PostBatchHandler(rep, logger))
	r.Get("/{id}", handlers.GetHandler(rep, logger))
	r.Get("/api/user/urls", handlers.GetAllUserURLs(rep, logger))
	r.Delete("/api/user/urls", handlers.DeleteUserURLs(rep, logger))
	r.Get("/ping", handlers.PingDataBase(rep, logger))

	logger.Info("server running")
	err := http.ListenAndServe(cfg.ServerAddress, r)
	if err != nil {
		logger.Fatal("ListenAndServe: ", err)
	}

	return nil
}
