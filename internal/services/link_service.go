package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"gorm.io/gorm"

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository"
)

// Jeu de caractères pour la génération des codes courts.
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// LinkService fournit la logique métier pour les liens.
type LinkService struct {
	linkRepo  repository.LinkRepository
	clickRepo repository.ClickRepository
}

// NewLinkService crée une nouvelle instance de LinkService.
func NewLinkService(linkRepo repository.LinkRepository, clickRepo repository.ClickRepository) *LinkService {
	return &LinkService{
		linkRepo:  linkRepo,
		clickRepo: clickRepo,
	}
}

// GenerateShortCode génère un code court aléatoire de la longueur spécifiée.
func (s *LinkService) GenerateShortCode(length int) (string, error) {
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}
	return string(b), nil
}

// CreateLink crée un nouveau lien raccourci avec un code unique.
func (s *LinkService) CreateLink(longURL string) (*models.Link, error) {

	var shortCode string
	const maxRetries = 5

	for i := 0; i < maxRetries; i++ {
		code, err := s.GenerateShortCode(6)
		if err != nil {
			return nil, fmt.Errorf("failed to generate short code: %w", err)
		}

		// Vérifie l'unicité du code généré.
		_, err = s.linkRepo.GetLinkByShortCode(code)

		if err != nil {
			// Le code est unique si 'record not found'.
			if errors.Is(err, gorm.ErrRecordNotFound) {
				shortCode = code
				break
			}
			return nil, fmt.Errorf("database error checking short code uniqueness: %w", err)
		}

		log.Printf("Short code '%s' already exists, retrying generation (%d/%d)...", code, i+1, maxRetries)

	}

	if shortCode == "" {
		return nil, errors.New("failed to generate unique short code after maximum retries")
	}

	link := &models.Link{
		ShortCode: shortCode,
		LongURL:   longURL,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.linkRepo.CreateLink(link)
	if err != nil {
		return nil, fmt.Errorf("failed to save link to database: %w", err)
	}

	return link, nil
}

// GetLinkByShortCode récupère un lien par son code court.
func (s *LinkService) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	return s.linkRepo.GetLinkByShortCode(shortCode)
}

// GetLinkStats récupère les statistiques pour un lien donné.
func (s *LinkService) GetLinkStats(shortCode string) (*models.Link, int, error) {
	// Récupère le lien.
	link, err := s.linkRepo.GetLinkByShortCode(shortCode)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get link: %w", err)
	}
	// Compte les clics.
	clickCount, err := s.clickRepo.CountClicksByLinkID(link.ID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count clicks: %w", err)
	}
	return link, clickCount, nil
}
