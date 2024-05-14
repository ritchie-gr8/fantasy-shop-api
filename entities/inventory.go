package entities

import "time"

type Inventory struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;"`
	PlayerID  string    `gorm:"type:varchar(64);not null;"`
	ItemID    uint64    `gorm:"type:bigint;not null;"`
	IsArchive bool      `gorm:"not null;default:false;"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime;"`
}
