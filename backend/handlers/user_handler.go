package handlers

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "pemm/models"
    _ "pemm/database"
	"golang.org/x/crypto/bcrypt"
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