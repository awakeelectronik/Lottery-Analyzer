package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"lottery-analyzer/api/middleware"
	"lottery-analyzer/internal/config"
	"lottery-analyzer/internal/repository"
	"lottery-analyzer/internal/service"
	"lottery-analyzer/pkg/database"

	_ "github.com/joho/godotenv/autoload"
)

type API struct {
	processor service.ProcessorService
}

func (api *API) processAnalysis(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	analysis, err := api.processor.ProcessAnalysis(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   analysis,
	})
}

func (api *API) BestNumbers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	ctx := r.Context()
	numbers, scores, err := api.processor.BestNumbers(ctx, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"numbers": numbers,
			"scores":  scores,
			"count":   len(numbers),
		},
	})
}

func (api *API) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func main() {
	cfg := config.Load()

	db, err := database.NewMySQL(cfg.Database.DSN)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	resultRepo := repository.NewResultRepository(db)
	scrapperService := service.NewScrapperService(resultRepo)
	processorService := service.NewProcessorService(scrapperService, resultRepo)

	api := &API{processor: processorService}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", api.healthCheck)
	mux.HandleFunc("/api/v1/analysis/process", api.processAnalysis)
	mux.HandleFunc("/api/v1/analysis/best-numbers", api.BestNumbers)

	handler := middleware.CORS(middleware.Logging(middleware.Recovery(mux)))

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("API server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown failed:", err)
	}
	log.Println("Server stopped")
}
