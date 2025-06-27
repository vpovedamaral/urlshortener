package monitor

import (
	"log"
	"net/http"
	"sync"
	"time"

	_ "github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository"
)

// UrlMonitor gère la surveillance périodique des URLs longues.
type UrlMonitor struct {
	linkRepo    repository.LinkRepository
	interval    time.Duration
	knownStates map[uint]bool // État connu de chaque URL: map[LinkID]estAccessible
	mu          sync.Mutex
}

// NewUrlMonitor crée une nouvelle instance de UrlMonitor.
func NewUrlMonitor(linkRepo repository.LinkRepository, interval time.Duration) *UrlMonitor {
	return &UrlMonitor{
		linkRepo:    linkRepo,
		interval:    interval,
		knownStates: make(map[uint]bool),
		mu:          sync.Mutex{},
	}
}

// Start lance la surveillance périodique des URLs dans une goroutine.
func (m *UrlMonitor) Start() {
	log.Printf("[MONITOR] Démarrage du moniteur d'URLs avec un intervalle de %v...", m.interval)
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	// Exécute une première vérification immédiatement.
	m.checkUrls()

	// Boucle principale déclenchée par le ticker.
	for range ticker.C {
		m.checkUrls()
	}
}

// checkUrls vérifie l'état de toutes les URLs longues enregistrées.
func (m *UrlMonitor) checkUrls() {
	log.Println("[MONITOR] Lancement de la vérification de l'état des URLs...")

	// Récupère toutes les URLs longues actives.
	links, err := m.linkRepo.GetAllLinks()
	if err != nil {
		log.Printf("[MONITOR] erruer lors de la récup des liens pour la surveillance : %v", err)
		return
	}

	for _, link := range links {
		// Vérifie l'accessibilité du lien.
		currentState := m.isUrlAccessible(link.LongURL)

		// Protège l'accès à la map knownStates.
		m.mu.Lock()
		previousState, exists := m.knownStates[link.ID]
		m.knownStates[link.ID] = currentState
		m.mu.Unlock()

		// Initialise l'état sans notifier si c'est la première vérification.
		if !exists {
			log.Printf("[MONITOR] État initial pour le lien %s (%s) : %s",
				link.ShortCode, link.LongURL, formatState(currentState))
			continue
		}

		// Notifie si l'état a changé.
		if currentState != previousState {
			log.Printf("[NOTIFICATION] Le lien %s (%s) est passé de %s à %s !",
				link.ShortCode, link.LongURL, formatState(previousState), formatState(currentState))
		}
	}
	log.Println("[MONITOR] Vérification de l'état des URLs terminée.")
}

// isUrlAccessible vérifie l'accessibilité d'une URL via une requête HTTP HEAD.
func (m *UrlMonitor) isUrlAccessible(url string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Effectue une requête HEAD.
	resp, err := client.Head(url)
	if err != nil {
		log.Printf("[MONITOR] Erreur d'accès à l'URL '%s': %v", url, err)
		return false
	}
	defer resp.Body.Close()

	// Codes 2xx ou 3xx indiquent une URL accessible.
	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

// formatState rend l'état plus lisible dans les logs.
func formatState(accessible bool) string {
	if accessible {
		return "ACCESSIBLE"
	}
	return "INACCESSIBLE"
}
