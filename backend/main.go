package main

import (
	"net/http"
	"html/template"
	"github.com/labstack/echo/v4"
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
	// ルート定義(/index)
	e.GET("/index", func(c echo.Context) error {
		return c.Render(http.StatusOK, "dogcreate_post.html", nil)
	})
	// サーバー起動
	e.Logger.Fatal(e.Start(":8080"))
}
