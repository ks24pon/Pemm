package handlers

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "pemm/models"
    _ "pemm/database"
	"golang.org/x/crypto/bcrypt"
	"github.com/labstack/echo-contrib/session" 
	"gorm.io/gorm"
)

// ユーザー関連の処理
type UserHandler struct {
	DB *gorm.DB
}

// パスワードをハッシュ化する関数
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// エラーハンドリング
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ユーザー登録処理
func (h *UserHandler) UserRegister(c echo.Context) error {
	// フォームデータを取得
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	// パスワードハッシュ化
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return c.String(http.StatusInternalServerError, "パスワードのハッシュ化に失敗しました")
	}

// ユーザーオブジェクト作成
user := models.User {
	Username: username,
	Email:		email,
	Password: hashedPassword,
}


// データーベースに保存
if err := h.DB.Create(&user).Error; err != nil {
	return c.String(http.StatusInternalServerError, "ユーザー登録に失敗しました")
}
// 成功後のレスポンス
return c.Redirect(http.StatusSeeOther, "/new")
}

// ログイン処理
func (h *UserHandler) Login(c echo.Context) error {
	// フォームからメールアドレスとパスワードを取得
	email := c.FormValue("email")
	password := c.FormValue("password")

	// ユーザー情報を格納変数
	var user models.User

	// メールアドレスでユーザー検索
	if err := h.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return c.Render(http.StatusUnauthorized, "login.html", map[string]interface{}{
		"csrf":	c.Get("csrf").(string),
		"message": "メールアドレスまたはパスワードが間違ってます",
		})
	}

	// パスワード検証
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.Render(http.StatusUnauthorized, "login.html", map[string]interface{}{
			"csrf": c.Get("csrf").(string),
			"message": "メールアドレスまたはパスワードが間違ってます",
		})
	}

	// セッションの作成
	sess, err := session.Get("session", c)

	if err != nil {
		return c.String(http.StatusInternalServerError, "セッションの取得に失敗しました")
	}

	// セッションににユーザーを保存
	sess.Values["user_id"] = user.ID
	sess.Values["username"] = user.Username

	// セッションの保存
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return c.String(http.StatusInternalServerError, "セッションの保存に失敗しました")
	}

	// 成功後リダイレクト(投稿画面)
	return c.Redirect(http.StatusSeeOther, "/new")
}