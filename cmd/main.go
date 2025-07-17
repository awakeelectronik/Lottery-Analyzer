package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lottery-analyzer/internal/config"
	"lottery-analyzer/internal/repository"
	"lottery-analyzer/internal/service"
	"lottery-analyzer/pkg/database"

	_ "github.com/joho/godotenv/autoload"
)

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	log.Println("Starting lottery analysis...")
	start := time.Now()

	analysis, err := processorService.ProcessAnalysis(ctx)
	if err != nil {
		log.Fatal("Analysis failed:", err)
	}

	log.Printf("Analysis completed in %v", time.Since(start))
	log.Printf("Best numbers: %v", analysis.BestNumbers[:10])
	log.Printf("Total processed: %d", analysis.TotalProcessed)
	log.Printf("Days analyzed: %d", analysis.DaysAnalyzed)
}
