package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

func main() {
	// 新しいEchoインスタンス生成
	e := echo.New()
	// ルートエンドポイントを定義
	e.GET("/", func(c echo.Context) error {
		// クライアントに返してレスポンスを返す
		return c.String(http.StatusOK, "Hello,Pemm！")
	})
	// サーバー起動
	e.Logger.Fatal(e.Start(":8080"))
	// テスト
}