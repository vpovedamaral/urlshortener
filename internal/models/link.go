package models

import (
	"time"

	"gorm.io/gorm"
)

// Link représente un lien raccourci dans la base de données.
type Link struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ShortCode string         `gorm:"unique;not null;size:6;index" json:"short_code"`
	LongURL   string         `gorm:"not null;type:text" json:"long_url"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relation un-à-plusieurs avec les clics.
	Clicks []Click `gorm:"foreignKey:LinkID" json:"clicks,omitempty"`
}

func (Link) TableName() string {
	return "links"
}
