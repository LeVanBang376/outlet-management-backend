package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"magnolia-test-backend/internal/evidence"
	middleware "magnolia-test-backend/internal/middlewares"
	"magnolia-test-backend/internal/outlet"
	"magnolia-test-backend/internal/sales"
	"magnolia-test-backend/internal/working_schedule"
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

	// Outlet
	outletRepository := outlet.NewRepository(db)
	outletService := outlet.NewService(outletRepository)
	outletHandler := outlet.NewHandler(outletService)

	// Working Schedule
	workingScheduleRepository := working_schedule.NewRepository(db)
	workingScheduleService := working_schedule.NewService(workingScheduleRepository)
	workingScheduleHandler := working_schedule.NewHandler(workingScheduleService)

	// Evidence
	evidenceRepository := evidence.NewRepository(db)
	evidenceService := evidence.NewService(evidenceRepository)
	evidenceHandler := evidence.NewHandler(evidenceService)

	// Sales
	salesRepository := sales.NewRepository(db)
	salesService := sales.NewService(salesRepository)
	salesHandler := sales.NewHandler(salesService)

	// Apply middlewares - Recovery and CORS
	handler := middleware.CORS([]string{"http://localhost:3000"})(mux)

	outlet.RegisterRoutes(mux, outletHandler)
	working_schedule.RegisterRoutes(mux, workingScheduleHandler)
	evidence.RegisterRoutes(mux, evidenceHandler)
	sales.RegisterRoutes(mux, salesHandler)

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
