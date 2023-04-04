package models

import (
	"gorm.io/gorm"
	"time"
)

type ProtocolDB struct {
	ID            uint           `gorm:"primary_key"`
	ProtocolNumber string         `gorm:"unique;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (P *ProtocolDB) TableName() string {
	return "protocolos"
}
