package models

import "time"

// Click représente un événement de clic sur un lien raccourci.
type Click struct {
	ID        uint `gorm:"primaryKey"`
	LinkID    uint `gorm:"index"`
	Link      Link `gorm:"foreignKey:LinkID"`
	Timestamp time.Time
	UserAgent string `gorm:"size:255"`
	IPAddress string `gorm:"size:50"`
}

func (Click) TableName() string {
	return "clicks"
}

type ClickEvent struct {
	LinkID    uint      `json:"link_id"`
	Timestamp time.Time `json:"timestamp"`
	UserAgent string    `json:"user_agent"`
	IPAddress string    `json:"ip_address"`
}
