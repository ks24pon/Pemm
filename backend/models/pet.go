package models

import (
	"time"
)

// 犬の登録構造体
type Dog struct {
	// 犬のID
    ID uint `json:"id™`
	// ユーザーID
	UseID uint `json: ™use_id™`
	// 名前(必須)
	Name string `json: ™name™ validate:required"`
	// ニックネーム
	Nickname string `json:"nickname" validate:required validate:"max=50"`
	// 種別
	Breed string `json: "breed" validate:required`
	// 性別
	Gender string `json: "gender" validate:required`
	// 誕生日
	Birthday *time.Time `json: "birthday" validate:required`
	// 登録日時
	CreatedAt time.Time `json: "created_at"`











}