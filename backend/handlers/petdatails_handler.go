package handlers

import (
	"net/http"
	"pemm/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// ペット詳細専用ハンドラ
type PetDetailHandler struct {
	DB *gorm.DB
}

// ペット詳細登録処理
func (h *PetDetailHandler) PetDetail(c echo.Context) error {
	// フォームデータ取得
	name := c.FormValue("name")
	type := c.FormValue("type")
	breed := c.FormValue("breed")
	gender := c.FormValue("gender")
	age := c.FormValue("age")
}