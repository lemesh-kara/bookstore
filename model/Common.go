package model

import "time"

type Common struct {
	ID        uint64     `gorm:"primary_key;autoIncrement" json:"id,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
}
