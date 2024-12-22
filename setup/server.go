package setup

import (
	"archive_lib/handler"
	"archive_lib/repo"
	"archive_lib/usecase"
	"archive_lib/util"
	"archive_lib/util/logger"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

const Timeout = 5

func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	db, err := ConnectDB()
	if err != nil {
		log.Fatalf("unable to connect to the database: %v", err)
	}
	defer db.Close()

	bcrypt := util.NewBcrypt()
	jwt := util.NewJWT()
	logger.SetLogger(logger.NewLogrusLogger())
	util.FormatValidatedField()

	txRepo := repo.NewTransactionRepo(db)

	userRepo := repo.NewUserRepo(db)
	userUsecase := usecase.NewUserUsecase(userRepo, bcrypt, jwt)
	userHandler := handler.NewUserHandler(userUsecase)

	bookRepo := repo.NewBookRepo(db)
	bookUsecase := usecase.NewBookUsecase(bookRepo)
	bookHandler := handler.NewBookHandler(bookUsecase)

	borrowRepo := repo.NewBorrowRepo(db)
	borrowUsecase := usecase.NewBorrowUsecase(borrowRepo, bookRepo, txRepo)
	borrowHandler := handler.NewBorrowHandler(borrowUsecase)

	handlers := NewHandlers(&userHandler, &bookHandler, &borrowHandler)
	router := NewRouter(handlers)

	s := &http.Server{
		Addr:         ":" + os.Getenv("SERVER_PORT"),
		Handler:      router,
		ReadTimeout:  Timeout * time.Second,
		WriteTimeout: Timeout * time.Second,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	<-ctx.Done()

	log.Println("Server exited gracefully")
}
