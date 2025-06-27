package repository

import (
	"github.com/axellelanca/urlshortener/internal/models"
	"gorm.io/gorm"
)

// ClickRepository définit les méthodes d'accès aux données pour les clics.
type ClickRepository interface {
	CreateClick(click *models.Click) error
	CountClicksByLinkID(linkID uint) (int, error)
}

// GormClickRepository implémente ClickRepository avec GORM.
type GormClickRepository struct {
	db *gorm.DB
}

// NewClickRepository crée une nouvelle instance de GormClickRepository.
func NewClickRepository(db *gorm.DB) *GormClickRepository {
	return &GormClickRepository{db: db}
}

// CreateClick insère un nouvel enregistrement de clic.
func (r *GormClickRepository) CreateClick(click *models.Click) error {
	return r.db.Create(click).Error
}

// CountClicksByLinkID compte le nombre total de clics pour un lien donné.
func (r *GormClickRepository) CountClicksByLinkID(linkID uint) (int, error) {
	var count int64
	err := r.db.Model(&models.Click{}).Where("link_id = ?", linkID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
