package models

import (
	"time"
	"gorm.io/gorm"
)

// 犬の投稿情報の構造体
type Post struct {
	// 一意の識別子
    ID uint `json:"id" gorm:"primary_key"`
	//犬の名前(最大255文字、空にできない)
	Name string `json:"name" gorm:"size:255;not null"`
	// 投稿に関する説明（Text型）
	Description string `json:"description" gorm:"type:text"`
	//写真のファイルパスを保存(最大255文字のファイパス)
	PhotoPath string `json:"photo_path" gorm:"gorm:size:255"`
	// GORMで自動的に作成日時を保存
	CreatedAt time.Time `json:"created_at"`
	// GORMで自動的に更新日時を保存
	UpdatedAt time.Time `json:"updated_at"`
	// ソフトデリート(論理削除)
	// gorm:"index"によりデーターベースのインデックスとしても扱う
	DeletedAt gorm.DeletedAt `gorm:"index"`
}