package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"pemm/models"
)
// データーベースで使用
var DB *gorm.DB

// データーベースの初期化を行う関数
func InitDB() {
	// 環境変数からデーターベース接続情報を取得
	dsn := os.Getenv("DB_DSN")
	// エラーハンドリング用の変数宣言
	var err error
	// GORMを通じてMysqlに接続して接続に成功するとDB変数に格納
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// エラーハンドリング(接続が失敗した場合)
	if err != nil {
		// log.Fatalでエラーメッセージと共にerrをの内容を表示
		log.Fatal("データーベース接続に失敗いたしました", err)
	}
	// 成功後
	log.Println("データーベース接続に成功いたしました")

	// 自動マイグレーション(テーブルを自動で作成)
	if err := DB.AutoMigrate(&models.Post{}); err != nil {
		log.Fatal("テーブルの自動作成に失敗しました:", err)
	} else {
		log.Println("テーブルのマイグレーションが成功しました")
	}
}