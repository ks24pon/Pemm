package handlers

import (
	"net/http"
	"pemm/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)


// ニックネームの処理
type PetHandler struct {
	DB *gorm.DB
}

// ペット登録(ニックネーム)処理
func (h *PetHandler) PetRegister(c echo.Context) error {
	// フォームデータ取得
	nickname := c.FormValue("nickname")

	// ペットオブジェクト作成
	pet := models.Pet {
		Nickname: nickname,
	}

	// データーベースに保存
	if err := h.DB.Create(&pet).Error; err != nil {
		return c.String(http.StatusInternalServerError, "ニックネーム登録に失敗しました")
	}

	// 成功後のレスポンス
	return c.Redirect(http.StatusSeeOther, "/new")
}