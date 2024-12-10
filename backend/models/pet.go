package models

import (
	"time"

	_ "gorm.io/gorm"
)

type Pet struct {
	// ペットID
	ID uint `json:"id"`
	// ユーザーID
	// UserID int    `json:"user_id"`
	// ニックネーム(3~50文字)
	Nickname string `json:"nickname" validate:"required,min=3,max=50"`
	// 作成日時
	CreatedAt time.Time `json:"created_at"`
}
