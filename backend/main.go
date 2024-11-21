package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"pemm/database"
	"pemm/handlers"
	_ "github.com/gorilla/csrf"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// HTMLをレンタリングする構造体
type HTMLTemplateRender struct {
	templates *template.Template
}

// Echoとtemplateパッケージを使ってレンダリング
func (render *HTMLTemplateRender) Render(writer io.Writer, name string, data interface{}, c echo.Context) error {
	return render.templates.ExecuteTemplate(writer, name, data)
}

func main() {
	// データベースの初期化
	database.InitDB()
	// 新しいEchoインスタンス生成
	e := echo.New()
	// UserHandlerのインスタンス生成
	UserHandler := &handlers.UserHandler {
		DB: database.DB,
	}
	// テンプレートの設定
	render := &HTMLTemplateRender{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = render
	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CSRF対策
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		// CSRFトークンを保存するクッキーの名前
		CookieName: "csrf_token",
		// フォームでトークンを受け取る方法を指定
		TokenLookup: "form:csrf_token",
		// クッキーをHTTPS通信のみで崇信するか
		CookieSecure: false,
		// クッキーをJavaScriptからアクセス不可にするかどうか
		CookieHTTPOnly: true,
	}))
	
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
	// 編集画面
	e.GET("/posts/:id/edit", handlers.EditPost)
	// 編集処理
	e.POST("/posts/:id/update", handlers.UpdatePost)
	// 詳細画面
	e.GET("/posts/:id", handlers.ShowPost)
	// 削除処理
	e.POST("/posts/:id/delete", handlers.DeletePost)
	// ユーザー登録画面(/register)
	e.GET("/register", func(c echo.Context) error {
		data := map[string]interface{}{
			"csrf": c.Get("csrf").(string),
		}
		return c.Render(http.StatusOK, "user_register.html", data)
	})
	// ユーザー登録処理
	e.POST("/register", UserHandler.UserRegister)
	// ルート一覧をターミナルに出力
	for _, route := range e.Routes() {
		log.Printf("Method: %s, Path: %s, Name: %s\n", route.Method, route.Path, route.Name)
	}
	e.File("/favicon.ico", "favicon.ico")
	// サーバー起動
	e.Logger.Fatal(e.Start(":8080"))
}
