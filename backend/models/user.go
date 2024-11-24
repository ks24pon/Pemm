package models

import (
	"time"
	_ "gorm.io/gorm"
)

type User struct {
	// 一意の識別子
	ID uint `json:"id" gorm:"primary_key"`
	// ユーザーの名前
	Username string `json:"name" gorm:"size:255;not null" validate:"min=8"`
	// メールアドレス
	Email string `json:"email" gorm:"size:255;unique;not null"`
	// パスワード
	Password string `json:"password" gorm:"size:255;not null"`
	// GORMで自動的に作成日時を保存
	CreatedAt time.Time `json:"created_at"`
}