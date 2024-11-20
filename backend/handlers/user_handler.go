package handlers

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "pemm/models"
    _ "pemm/database"
	"gorm.io/gorm"
)

// ユーザー関連の処理
type UserHandler struct {
	DB *gorm.DB
}

// ユーザー登録処理
func (h *UserHandler) UserRegister(c echo.Context) error {
	// フォームデータを取得
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")


// ユーザーオブジェクト作成
user := models.User {
	Username: username,
	Email:		email,
	Password: password,
}

// データーベースに保存
if err := h.DB.Create(&user).Error; err != nil {
	return c.String(http.StatusInternalServerError, "ユーザー登録に失敗しました")
}
// 成功後のレスポンス
return c.Redirect(http.StatusSeeOther, "/new")
}