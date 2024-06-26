package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shafaalafghany/segokuning-social-app/config"
	"github.com/shafaalafghany/segokuning-social-app/internal/common/utils/validation"
	commentHandler "github.com/shafaalafghany/segokuning-social-app/internal/handler/comment"
	friendHandler "github.com/shafaalafghany/segokuning-social-app/internal/handler/friend"
	imageHandler "github.com/shafaalafghany/segokuning-social-app/internal/handler/image"
	postHandler "github.com/shafaalafghany/segokuning-social-app/internal/handler/post"
	userHandler "github.com/shafaalafghany/segokuning-social-app/internal/handler/user"
	"github.com/shafaalafghany/segokuning-social-app/internal/repository"
	"github.com/shafaalafghany/segokuning-social-app/pkg/db"
	"github.com/shafaalafghany/segokuning-social-app/pkg/logger"
)

func Run(cfg *config.Configuration) {
	var validate *validator.Validate

	pgx := db.NewPsqlDB(cfg)

	validate = validator.New()
	if err := validation.RegisterCustomValidation(validate); err != nil {
		log.Fatalf("error register custom validation")
	}

	logger, err := logger.Initialize(*cfg)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	r := chi.NewRouter()

	ur := repository.NewUserRepo(pgx, logger)
	fr := repository.NewFriendRepo(pgx, logger)
	cr := repository.NewCommentRepo(pgx, logger)
	pr := repository.NewPostRepo(pgx, logger)

	r.Handle("/metrics", promhttp.Handler())
	r.Route("/v1", func(r chi.Router) {
		// r.Use(promotheus.PrometheusMiddleware)
		userHandler.NewUserHandler(r, ur, validate, *cfg, logger)
		friendHandler.NewFriendHandler(r, ur, fr, validate, *cfg, logger)
		postHandler.NewPostHandler(r, ur, pr, validate, *cfg, logger)
		commentHandler.NewCommentHandler(r, fr, cr, pr, validate, *cfg, logger)
		imageHandler.NewImageHandler(r, *validate, *cfg, logger)
	})

	s := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}
	go func() {
		fmt.Println("Listen and Serve at port 8080")
		if err := s.ListenAndServe(); err != nil {
			log.Fatalf("error in ListenAndServe: %s", err)
		}
	}()
	log.Print("Server Started")

	stopped := make(chan os.Signal, 1)
	signal.Notify(stopped, os.Interrupt)
	<-stopped

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("shutting down gracefully...")
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("error in Server Shutdown: %s", err)
	}
	fmt.Println("server stopped")
}
