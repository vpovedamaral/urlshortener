package services

import (
	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository"
)

// ClickService fournit la logique métier pour les clics.
type ClickService struct {
	clickRepo repository.ClickRepository
}

// NewClickService crée une nouvelle instance de ClickService.
func NewClickService(clickRepo repository.ClickRepository) *ClickService {
	return &ClickService{
		clickRepo: clickRepo,
	}
}

// RecordClick enregistre un nouvel événement de clic.
func (s *ClickService) RecordClick(click *models.Click) error {
	return s.clickRepo.CreateClick(click)
}

// GetClicksCountByLinkID récupère le nombre total de clics pour un lien donné.
func (s *ClickService) GetClicksCountByLinkID(linkID uint) (int, error) {
	return s.clickRepo.CountClicksByLinkID(linkID)
}
