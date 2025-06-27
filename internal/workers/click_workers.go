package workers

import (
	"log"

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository"
)

// StartClickWorkers lance un pool de workers pour traiter les événements de clic.
func StartClickWorkers(workerCount int, clickEventsChan <-chan models.ClickEvent, clickRepo repository.ClickRepository) {
	log.Printf("Starting %d click worker(s)...", workerCount)
	for i := 0; i < workerCount; i++ {
		go clickWorker(clickEventsChan, clickRepo)
	}
}

// clickWorker traite les événements de clic de manière asynchrone.
func clickWorker(clickEventsChan <-chan models.ClickEvent, clickRepo repository.ClickRepository) {
	for event := range clickEventsChan {
		// Convertit l'événement en modèle Click.
		click := &models.Click{
			LinkID:    event.LinkID,
			Timestamp: event.Timestamp,
			UserAgent: event.UserAgent,
			IPAddress: event.IPAddress,
		}

		// Persiste le clic en base de données.
		err := clickRepo.CreateClick(click)
		if err != nil {
			log.Printf("ERROR: Failed to save click for LinkID %d: %v",
				event.LinkID, err)
		} else {
			log.Printf("Click recorded successfully for LinkID %d", event.LinkID)
		}
	}
}
