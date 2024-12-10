package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"pemm/database"
	"pemm/handlers"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// HTMLをレンダリングする構造体
type HTMLTemplateRender struct {
	templates *template.Template
}

// Echoとtemplateパッケージを使ってレンダリング
func (render *HTMLTemplateRender) Render(writer io.Writer, name string, data interface{}, c echo.Context) error {
	return render.templates.ExecuteTemplate(writer, name, data)
}

// カスタムバリデーション構造体
type CustomValidation struct {
	validator *validator.Validate
}

// カスタムバリデーション関数
func (cv *CustomValidation) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// バリデーション関数(8文字以上)
func CheckMinLength(fl validator.FieldLevel) bool {
	return len(fl.Field().String()) >= 8
}

func main() {
	// データベースの初期化
	database.InitDB()

	// 新しいEchoインスタンス生成
	e := echo.New()

	// UserHandlerのインスタンス生成
	UserHandler := &handlers.UserHandler{
		DB: database.DB,
	}

	// ニックネームインスタンス作成
	PetHandler := &handlers.PetHandler{
		DB: database.DB,
	}

	// テンプレートの設定
	render := &HTMLTemplateRender{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = render

	// バリデーションインスタンスの作成
	validate := validator.New()

	// カスタムバリデーション登録
	if err := validate.RegisterValidation("minlength8", CheckMinLength); err != nil {
		log.Fatalf("カスタムバリデーションの登録に失敗しました: %v", err)
	}

	// Echoにカスタムバリデーションを登録
	e.Validator = &CustomValidation{
		validator: validate,
	}

	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CSRF対策
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookieName:     "csrf_token",
		TokenLookup:    "form:csrf_token",
		ContextKey:     "csrf",
		CookieSecure:   false,
		CookieHTTPOnly: true,
		CookiePath:     "/",
	}))

	// Top画面のルート
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "top.html", nil)
	})
	// 静的ファイルの設定
	e.Static("/static", "static")

	// 静的ファイルの設定
	e.Static("/uploads", "uploads")

	// エラーハンドリング
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		c.Logger().Error("エラーが発生しました", err)
		c.String(http.StatusInternalServerError, "内部エラーが発生しました")
	}

	// 投稿画面ルート(/new)
	e.GET("/new", func(c echo.Context) error {
		// CSRFトークンをテンプレートに渡す
		data := map[string]interface{}{
			"csrf": c.Get("csrf").(string),
		}
		return c.Render(http.StatusOK, "dogcreate_post.html", data)
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

	// ニックネーム登録画面(nickname)
	e.GET("/nickname", func(c echo.Context) error {
		data := map[string]interface{}{
			"csrf": c.Get("csrf").(string),
		}
		return c.Render(http.StatusOK, "nickname_register.html", data)
	})

	// ニックネーム登録処理(/nickname)
	e.POST("/nickname", PetHandler.PetRegister)

	// ペット登録画面(/pet/register)
	e.GET("/pet/register",func(c echo.Context) error {
		data := map[string]interface{}{
			"csrf": c.Get("csrf").(string),
		}
		return c.Render(http.StatusOK, "pet_register.html", data)
	})

	// ペット登録処理

	// ルート一覧をターミナルに出力
	for _, route := range e.Routes() {
		log.Printf("Method: %s, Path: %s, Name: %s\n", route.Method, route.Path, route.Name)
	}

	e.File("/favicon.ico", "favicon.ico")

	// サーバー起動
	e.Logger.Fatal(e.Start(":8080"))
}
