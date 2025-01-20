package models

import (
	"time"
	_ "gorm.io/gorm"
)

// ペット詳細Modelの定義
type PetDetail struct {
	// ペットID
	ID uint `gorm:"primaryKey"`
	// ペット名前
	Name string `gorm:"size:255;not null" validate:"required"`
	// ペット種類
	PetType string `gorm:"size:255;not null" validate:"required"`
	// ペット種別
	Breed string `gorm:"size:255;not null" validate:"required"`
	// 性別
	Gender string `gorm:"size:10;not null" validate:"required,oneof=male female"`
	// 年齢
	Age int `gorm:"not null" validate:"gte=0"`
	// 日時作成
	CreatedAt time.Time
	// 日時更新
	UpdatedAt time.Time
}