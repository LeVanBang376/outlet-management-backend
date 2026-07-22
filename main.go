package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"magnolia-test-backend/internal/handler"
	middleware "magnolia-test-backend/internal/middlewares"
	"magnolia-test-backend/internal/repository"
	"magnolia-test-backend/internal/routes"
	"magnolia-test-backend/internal/service"
	"magnolia-test-backend/internal/worker"
	"net"
	"net/http"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

const (
	_shutdownPeriod      = 15 * time.Second
	_shutdownHardPeriod  = 3 * time.Second
	_readinessDrainDelay = 5 * time.Second
)

var isShuttingDown atomic.Bool

func main() {
	// Setup signal context
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if isShuttingDown.Load() {
			http.Error(w, "Shutting down", http.StatusServiceUnavailable)
			return
		}
		fmt.Fprintln(w, "OK")
	})

	// Database
	db, err := sql.Open("sqlite", "app.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Optional: verify connection
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal(err)
	}

	// if err := goose.Reset(db, "migrations"); err != nil {
	// 	log.Fatal(err)
	// }

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal(err)
	}

	log.Println("Database migrated successfully.")

	mux := http.NewServeMux()

	// Repositories
	outletRepository := repository.NewOutletRepository(db)
	workingScheduleRepository := repository.NewWorkingScheduleRepository(db)
	evidenceRepository := repository.NewEvidenceRepository(db)
	salesRepository := repository.NewSalesRepository(db)
	fileRepository := repository.NewFileRepository(db)

	// Worker
	worker := worker.NewWorker(workingScheduleRepository)

	// Services
	outletService := service.NewOutletService(db, outletRepository, workingScheduleRepository, worker)
	workingScheduleService := service.NewWorkingScheduleService(db, workingScheduleRepository, outletRepository, evidenceRepository, worker)
	evidenceService := service.NewEvidenceService(db, evidenceRepository, fileRepository)
	salesService := service.NewSalesService(db, salesRepository)
	fileService := service.NewFileService(db, fileRepository)

	// Handlers
	outletHandler := handler.NewOutletHandler(outletService)
	workingScheduleHandler := handler.NewWorkingScheduleHandler(workingScheduleService)
	evidenceHandler := handler.NewEvidenceHandler(evidenceService)
	salesHandler := handler.NewSalesHandler(salesService)
	fileHandler := handler.NewFileHandler(fileService)

	// Apply middlewares - Recovery and CORS
	handler := middleware.CORS([]string{"http://localhost:3000"})(mux)

	routes.RegisterOutletRoutes(mux, outletHandler)
	routes.RegisterWorkingScheduleRoutes(mux, workingScheduleHandler)
	routes.RegisterEvidenceRoutes(mux, evidenceHandler)
	routes.RegisterSalesRoutes(mux, salesHandler)
	routes.RegisterFileRoutes(mux, fileHandler)

	// Ensure in-flight requests aren't cancelled immediately on SIGTERM
	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
		BaseContext: func(_ net.Listener) context.Context {
			return ongoingCtx
		},
	}

	go func() {
		log.Println("Server starting on :8080.")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Wait for signal
	<-rootCtx.Done()
	stop()
	isShuttingDown.Store(true)
	log.Println("Received shutdown signal, shutting down.")

	// Give time for readiness check to propagate
	time.Sleep(_readinessDrainDelay)
	log.Println("Readiness check propagated, now waiting for ongoing requests to finish.")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), _shutdownPeriod)
	defer cancel()
	err = server.Shutdown(shutdownCtx)
	stopOngoingGracefully()
	if err != nil {
		log.Println("Failed to wait for ongoing requests to finish, waiting for forced cancellation.")
		time.Sleep(_shutdownHardPeriod)
	}

	log.Println("Server shut down gracefully.")
}
