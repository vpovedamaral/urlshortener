package repository

import (
	"github.com/axellelanca/urlshortener/internal/models"
	"gorm.io/gorm"
)

// LinkRepository est une interface qui définit les méthodes d'accès aux données
// pour les opérations CRUD sur les liens.
type LinkRepository interface {
	CreateLink(link *models.Link) error
	GetLinkByShortCode(shortCode string) (*models.Link, error)
	GetAllLinks() ([]models.Link, error)
	CountClicksByLinkID(linkID uint) (int, error)
}

// TODO (fait):  GormLinkRepository est l'implémentation de LinkRepository utilisant GORM.
type GormLinkRepository struct {
	db *gorm.DB
}

// NewLinkRepository crée et retourne une nouvelle instance de GormLinkRepository.
// Cette fonction retourne *GormLinkRepository, qui implémente l'interface LinkRepository.
func NewLinkRepository(db *gorm.DB) *GormLinkRepository {
	return &GormLinkRepository{db: db}
}

// CreateLink insère un nouveau lien dans la base de données.
func (r *GormLinkRepository) CreateLink(link *models.Link) error {
	//  Utiliser GORM pour créer un nouvel enregistrement (link) dans la table des liens.
	return r.db.Create(link).Error
}

// GetLinkByShortCode récupère un lien de la base de données en utilisant son shortCode.
// Il renvoie gorm.ErrRecordNotFound si aucun lien n'est trouvé avec ce shortCode.
func (r *GormLinkRepository) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	var link models.Link
	// TODO 2: Utiliser GORM pour trouver un lien par son ShortCode.
	// La méthode First de GORM recherche le premier enregistrement correspondant et le mappe à 'link'.

	err := r.db.Where("short_code = ?", shortCode).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// GetAllLinks récupère tous les liens de la base de données.
// Cette méthode est utilisée par le moniteur d'URLs.
func (r *GormLinkRepository) GetAllLinks() ([]models.Link, error) {
	var links []models.Link
	// TODO 3: Utiliser GORM pour récupérer tous les liens.
	err := r.db.Find(&links).Error
	return links, err
}

// CountClicksByLinkID compte le nombre total de clics pour un ID de lien donné.
func (r *GormLinkRepository) CountClicksByLinkID(linkID uint) (int, error) {
	var count int64 // GORM retourne un int64 pour les comptes
	// TODO 4: Utiliser GORM pour compter les enregistrements dans la table 'clicks'
	// où 'LinkID' correspond à l'ID du lien donné.

	err := r.db.Model(&models.Click{}).Where("link_id = ?", linkID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
