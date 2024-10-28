package main

import (
	"net/http"
	"pemm/database"
	"html/template"
	"github.com/labstack/echo/v4"
	"pemm/handlers"
	"github.com/labstack/echo/v4/middleware"
	"io"
)

// HTMLをレンタリングする構造体
type HTMLTemplateRender struct {
	templates *template.Template
}
// Echoとtemplateパッケージを使ってレンダリング
func (render *HTMLTemplateRender) Render(writer io.Writer, name string, data interface{}, c echo.Context) error {
	return  render.templates.ExecuteTemplate(writer, name, data)
}



func main() {
	// データベースの初期化
	database.InitDB()
	// 新しいEchoインスタンス生成
	e := echo.New()
	// テンプレートの設定
	render := &HTMLTemplateRender {
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = render
	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// ルートを定義
	e.GET("/", func(c echo.Context) error {
		// クライアントに返してレスポンスを返す
		return c.String(http.StatusOK, "Hello,Pemm！")
	})
	// 静的ファイルの設定
	e.Static("/uploads", "uploads")
	// エラーハンドリング
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		// エラーが発生したことをログに記録
		c.Logger().Error("エラーが発生しました", err)
		// クライアント内部エラーが発生したら返す
		c.String(http.StatusInternalServerError, "内部エラーが発生しました")
	}
	// 投稿画面ルート(/new)
	e.GET("/new", func(c echo.Context) error {
		return c.Render(http.StatusOK, "dogcreate_post.html", nil)
	})
	// 投稿処理
	e.POST("/posts", handlers.CreatePost)
	// 投稿一覧
	e.GET("/index", handlers.ListPosts)
	// サーバー起動
	e.Logger.Fatal(e.Start(":8080"))

}
