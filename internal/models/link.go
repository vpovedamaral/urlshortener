package models

import (
	"time"

	"gorm.io/gorm"
)

// TODO : Créer la struct Link
// Link représente un lien raccourci dans la base de données.
// Les tags `gorm:"..."` définissent comment GORM doit mapper cette structure à une table SQL.
// ID qui est une primaryKey
// Shortcode : doit être unique, indexé pour des recherches rapide (voir doc), taille max 10 caractères
// LongURL : doit pas être null
// CreateAt : Horodatage de la créatino du lien

type Link struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ShortCode string         `gorm:"unique;not null;size:6;index" json:"short_code"`
	LongURL   string         `gorm:"not null;type:text" json:"long_url"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relation avec les clics - un lien peut avoir plusieurs clics
	Clicks []Click `gorm:"foreignKey:LinkID" json:"clicks,omitempty"`
}

func (Link) TableName() string {
	return "links"
}
